package sizer

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jmsilvadev/go-pack-optimizer/pkg/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// SizerInterface defines the methods that any Sizer implementation must provide.
type SizerInterface interface {
	GetAllSizes() ([]int, error)
	AddSize(size int) error
	RemoveSize(size int) error
	Close() error
}

// Sizer is responsible for interacting with pack sizes stored in a LevelDB database.
type Sizer struct {
	db     *leveldb.DB
	logger logger.Logger
}

// NewSizer opens or creates a LevelDB instance and populates it with default sizes if needed.
func NewSizer(path string, l logger.Logger) (*Sizer, error) {
	db, err := leveldb.OpenFile(path, &opt.Options{
		ErrorIfMissing: false,
	})
	if err != nil {
		return nil, err
	}

	sizer := &Sizer{
		db:     db,
		logger: l,
	}

	if err := sizer.Populate(); err != nil {
		return nil, err
	}

	return sizer, nil
}

// Populate fills the database with default pack sizes if it is empty or missing required data.
func (s *Sizer) Populate() error {
	data, err := s.db.Get([]byte("packs"), nil)
	if err == leveldb.ErrNotFound || len(data) == 0 {
		s.logger.Info("Table not found, creating and populating table...")

		defaultSizes := []int{250, 500, 1000, 2000, 5000}

		for _, size := range defaultSizes {
			if err := s.db.Put([]byte(fmt.Sprintf("size_%d", size)), []byte(fmt.Sprintf("%d", size)), nil); err != nil {
				return fmt.Errorf("failed to insert default size %d: %v", size, err)
			}
		}

		if err := s.db.Put([]byte("packs"), []byte("populated"), nil); err != nil {
			return fmt.Errorf("failed to set populated flag: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if packs exist: %v", err)
	}

	return nil
}

// Close closes the DB connection
func (s *Sizer) Close() error {
	return s.db.Close()
}

// GetAllSizes returns all sizes from LevelDB sorted in descending order
func (s *Sizer) GetAllSizes() ([]int, error) {
	iter := s.db.NewIterator(nil, nil)
	defer iter.Release()

	var sizes []int
	for iter.Next() {
		key := string(iter.Key())
		if strings.HasPrefix(key, "size_") {
			sizeStr := key[len("size_"):]
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				continue
			}
			sizes = append(sizes, size)
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	return sizes, iter.Error()
}

// AddSize adds a new pack size to the database
func (s *Sizer) AddSize(size int) error {
	key := fmt.Sprintf("size_%d", size)
	return s.db.Put([]byte(key), []byte(fmt.Sprintf("%d", size)), nil)
}

// RemoveSize deletes a pack size from the database
func (s *Sizer) RemoveSize(size int) error {
	key := fmt.Sprintf("size_%d", size)
	exists, err := s.db.Has([]byte(key), nil)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("pack size not found")
	}
	return s.db.Delete([]byte(key), nil)
}
