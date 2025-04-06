package parallel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOrdered(t *testing.T) {
	t.Parallel()

	t.Run("Works with single item", func(t *testing.T) {
		t.Parallel()
		iter := Ordered([]int{2}, func(input int) int {
			return input * input
		})

		expectedIndex := []int{0}
		expected := []int{4}

		actualIndex := make([]int, 0)
		actual := make([]int, 0)
		for i, v := range iter {
			actualIndex = append(actualIndex, i)
			actual = append(actual, v)
		}

		require.Equal(t, expectedIndex, actualIndex, "returned index")
		require.Equal(t, expected, actual, "returned value")
	})

	t.Run("Works with no items", func(t *testing.T) {
		t.Parallel()
		iter := Ordered([]int{}, func(input int) int {
			return input * input
		})

		expectedIndex := []int{}
		expected := []int{}

		actualIndex := make([]int, 0)
		actual := make([]int, 0)
		for i, v := range iter {
			actualIndex = append(actualIndex, i)
			actual = append(actual, v)
		}

		require.Equal(t, expectedIndex, actualIndex, "returned index")
		require.Equal(t, expected, actual, "returned value")
	})

	t.Run("Returns correct order", func(t *testing.T) {
		t.Parallel()
		iter := Ordered([]int{1, 2, 3, 4, 5}, func(input int) int {
			time.Sleep(time.Duration(input) * time.Millisecond * 20)
			return input * input
		})

		expectedIndex := []int{0, 1, 2, 3, 4}
		expected := []int{1, 4, 9, 16, 25}

		actualIndex := make([]int, 0)
		actual := make([]int, 0)
		for i, v := range iter {
			actualIndex = append(actualIndex, i)
			actual = append(actual, v)
		}

		require.Equal(t, expectedIndex, actualIndex, "returned index")
		require.Equal(t, expected, actual, "returned value")
	})
}

const benchmarkNumFiles = 5

func BenchmarkLoopedOrderedIO(b *testing.B) {
	docs := make([]*simulatedDoc, benchmarkNumFiles)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchmarkNumFiles; i++ {
			docs[i] = simulateReadYAML(b)
		}
	}
}

func BenchmarkParallelOrderedIO(b *testing.B) {
	docs := make([]*simulatedDoc, benchmarkNumFiles)
	inputs := make([]int, benchmarkNumFiles)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i, v := range Ordered(inputs, func(_ int) *simulatedDoc {
			return simulateReadYAML(b)
		}) {
			docs[i] = v
		}
	}
}

func BenchmarkLoopedOrderedIOFolded(b *testing.B) {
	merged := new(simulatedDoc)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchmarkNumFiles; i++ {
			doc := simulateReadYAML(b)
			merged = merged.Merge(doc)
		}
	}
}

func BenchmarkParallelOrderedIOFolded(b *testing.B) {
	inputs := make([]int, benchmarkNumFiles)
	merged := new(simulatedDoc)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, v := range Ordered(inputs, func(_ int) *simulatedDoc {
			return simulateReadYAML(b)
		}) {
			merged = merged.Merge(v)
		}
	}
}
