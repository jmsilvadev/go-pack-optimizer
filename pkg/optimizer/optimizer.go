package optimizer

import (
	"fmt"
	"sort"

	"github.com/jmsilvadev/go-pack-optimizer/pkg/logger"
	"github.com/jmsilvadev/go-pack-optimizer/pkg/sizer"
)

// memo stores previously computed optimization results to avoid redundant calculations.
var memo = make(map[int]*result)

//go:generate mockgen -source=optimizer.go -destination=../../internal/handler/mocks/mock_optimizer.go -package=mocks

// OptimizerInterface defines the behavior expected from any optimizer implementation.
type OptimizerInterface interface {
	Load() error
	Calculate(itemsOrdered int) *OptimizationResult
	GetAllSizes() ([]int, error)
	AddSize(size int) error
	RemoveSize(size int) error
}

// Optimizer provides methods for calculating optimal packaging solutions.
type Optimizer struct {
	sizer  sizer.SizerInterface
	logger logger.Logger
	sizes  []int
}

// OptimizationResult holds the final output of a packaging optimization.
type OptimizationResult struct {
	PacksUsed  []int `json:"packs_used"`
	TotalItems int   `json:"total_items"`
	TotalPacks int   `json:"total_packs"`
}

// result is an internal struct used during recursive computation of packaging.
type result struct {
	packs      []int
	totalItems int
}

// New creates a new Optimizer instance and preloads the available pack sizes.
func New(s sizer.SizerInterface, l logger.Logger) *Optimizer {
	opt := &Optimizer{
		sizer:  s,
		logger: l,
	}
	if err := opt.Load(); err != nil {
		l.Info(fmt.Sprintf("Error loading sizes: %v\n", err))
	}
	return opt
}

// Load retrieves and caches the available pack sizes from the sizer, sorted in descending order.
func (opt *Optimizer) Load() error {
	sizes, err := opt.sizer.GetAllSizes()
	if err != nil {
		return err
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	opt.sizes = sizes
	return nil
}

// Calculate returns the best combination of pack sizes for the given number of items.
func (opt *Optimizer) Calculate(itemsOrdered int) *OptimizationResult {
	if len(opt.sizes) == 0 || itemsOrdered <= 0 {
		return &OptimizationResult{}
	}

	res := opt.discoverPackages(itemsOrdered)
	if res == nil {
		return &OptimizationResult{}
	}

	return &OptimizationResult{
		PacksUsed:  res.packs,
		TotalItems: res.totalItems,
		TotalPacks: len(res.packs),
	}
}

// discoverPackages recursively determines the most efficient combination of pack sizes
// that covers at least the requested number of items, minimizing totalItems and pack count.
func (opt *Optimizer) discoverPackages(items int) *result {
	if res, ok := memo[items]; ok {
		return res
	}

	var best *result

	for _, size := range opt.sizes {
		newRemaining := items - size
		var candidate *result

		if newRemaining <= 0 {
			candidate = &result{
				packs:      []int{size},
				totalItems: size,
			}
		} else {
			next := opt.discoverPackages(newRemaining)
			if next == nil {
				continue
			}
			candidate = &result{
				packs:      append([]int{size}, next.packs...),
				totalItems: size + next.totalItems,
			}
		}

		best = chooseBetter(best, candidate)
	}

	memo[items] = best
	return best
}

// chooseBetter compares two packaging results and returns the more efficient one.
func chooseBetter(a, b *result) *result {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if b.totalItems < a.totalItems {
		return b
	}
	if b.totalItems == a.totalItems && len(b.packs) < len(a.packs) {
		return b
	}
	return a
}

// GetAllSizes returns all available pack sizes.
func (opt *Optimizer) GetAllSizes() ([]int, error) {
	return opt.sizer.GetAllSizes()
}

// AddSize adds a new pack size to the system.
func (opt *Optimizer) AddSize(size int) error {
	err := opt.sizer.AddSize(size)
	if err != nil {
		return err
	}

	return opt.reloadValues()
}

// RemoveSize deletes a pack size from the system.
func (opt *Optimizer) RemoveSize(size int) error {
	err := opt.sizer.RemoveSize(size)
	if err != nil {
		return err
	}

	return opt.reloadValues()
}

func (opt *Optimizer) reloadValues() error {
	memo = make(map[int]*result)
	return opt.Load()
}
