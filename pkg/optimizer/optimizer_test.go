package optimizer_test

import (
	"fmt"
	"testing"

	"github.com/jmsilvadev/go-pack-optimizer/pkg/logger"
	"github.com/jmsilvadev/go-pack-optimizer/pkg/optimizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
)

type MockSizer struct {
	mock.Mock
}

func (m *MockSizer) GetAllSizes() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockSizer) AddSize(size int) error {
	args := m.Called(size)
	return args.Error(0)
}

func (m *MockSizer) RemoveSize(size int) error {
	args := m.Called(size)
	return args.Error(0)
}

func (m *MockSizer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewOptimizer(t *testing.T) {
	mockSizer := new(MockSizer)
	mockSizer.On("GetAllSizes").Return([]int{500, 1000, 250}, nil).Maybe()

	opt := optimizer.New(mockSizer, logger.New(zapcore.DebugLevel))

	assert.NotNil(t, opt)
	v, _ := opt.GetAllSizes()
	assert.Equal(t, 3, len(v))
	mockSizer.AssertExpectations(t)
}

func TestLoad(t *testing.T) {
	mockSizer := new(MockSizer)
	mockSizer.On("GetAllSizes").Return([]int{1000, 250, 500}, nil).Maybe()

	opt := optimizer.New(mockSizer, logger.New(zapcore.DebugLevel))
	err := opt.Load()

	assert.Nil(t, err)
	v, _ := opt.GetAllSizes()
	assert.Equal(t, []int{1000, 500, 250}, v)
	mockSizer.AssertExpectations(t)
}

func TestCalculate(t *testing.T) {
	mockSizer := new(MockSizer)

	mockSizer.On("GetAllSizes").Return([]int{1000, 500, 250}, nil).Maybe()
	opt := optimizer.New(mockSizer, logger.New(zapcore.DebugLevel))

	tests := []struct {
		itemsOrdered int
		expected     *optimizer.OptimizationResult
	}{
		{
			itemsOrdered: 501,
			expected: &optimizer.OptimizationResult{
				PacksUsed:  []int{500, 250},
				TotalItems: 750,
				TotalPacks: 2,
			},
		},
		{
			itemsOrdered: 1500,
			expected: &optimizer.OptimizationResult{
				PacksUsed:  []int{1000, 500},
				TotalItems: 1500,
				TotalPacks: 2,
			},
		},
		{
			itemsOrdered: 0,
			expected: &optimizer.OptimizationResult{
				PacksUsed:  nil,
				TotalItems: 0,
				TotalPacks: 0,
			},
		},
		{
			itemsOrdered: -100,
			expected: &optimizer.OptimizationResult{
				PacksUsed:  nil,
				TotalItems: 0,
				TotalPacks: 0,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("itemsOrdered=%d", test.itemsOrdered), func(t *testing.T) {
			result := opt.Calculate(test.itemsOrdered)

			assert.Equal(t, test.expected.PacksUsed, result.PacksUsed)
			assert.Equal(t, test.expected.TotalItems, result.TotalItems)
			assert.Equal(t, test.expected.TotalPacks, result.TotalPacks)
		})
	}

	mockSizer.AssertExpectations(t)
}

func TestAddSize(t *testing.T) {
	mockSizer := new(MockSizer)
	mockSizer.On("GetAllSizes").Return([]int{1000, 500, 250}, nil)
	mockSizer.On("AddSize", 300).Return(nil)

	opt := optimizer.New(mockSizer, logger.New(zapcore.DebugLevel))

	err := opt.AddSize(300)
	assert.Nil(t, err)

	mockSizer.AssertExpectations(t)
}

func TestRemoveSize(t *testing.T) {
	mockSizer := new(MockSizer)
	mockSizer.On("GetAllSizes").Return([]int{1000, 500, 250}, nil)
	mockSizer.On("RemoveSize", 300).Return(nil)

	opt := optimizer.New(mockSizer, logger.New(zapcore.DebugLevel))

	err := opt.RemoveSize(300)
	assert.Nil(t, err)

	mockSizer.AssertExpectations(t)
}

func TestGetAllSizes(t *testing.T) {
	mockSizer := new(MockSizer)
	mockSizer.On("GetAllSizes").Return([]int{250, 500, 1000}, nil).Maybe()

	opt := optimizer.New(mockSizer, logger.New(zapcore.DebugLevel))

	sizes, err := opt.GetAllSizes()
	assert.Nil(t, err)
	assert.Equal(t, []int{1000, 500, 250}, sizes)

	mockSizer.AssertExpectations(t)
}

func TestCalculate_EmptySizes(t *testing.T) {
	mockSizer := new(MockSizer)
	mockSizer.On("GetAllSizes").Return([]int{}, nil).Maybe()

	opt := optimizer.New(mockSizer, logger.New(zapcore.DebugLevel))

	result := opt.Calculate(501)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.TotalPacks)
	assert.Nil(t, result.PacksUsed)

	mockSizer.AssertExpectations(t)
}

func TestCalculate_InvalidItemsOrdered(t *testing.T) {
	mockSizer := new(MockSizer)
	mockSizer.On("GetAllSizes").Return([]int{500, 1000, 250}, nil).Maybe()

	opt := optimizer.New(mockSizer, logger.New(zapcore.DebugLevel))

	result := opt.Calculate(-100)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.TotalPacks)
	assert.Nil(t, result.PacksUsed)

	result = opt.Calculate(0)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.TotalPacks)
	assert.Nil(t, result.PacksUsed)

	mockSizer.AssertExpectations(t)
}
