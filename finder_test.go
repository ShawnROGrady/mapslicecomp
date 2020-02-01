package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

/*
var allFinders = map[string]func() finder{
	"slice_finder": newSliceFinder,
	"map_finder":   newMapFinder,
}
*/

var allFinders = []struct {
	name  string
	newFn func() finder
}{
	{
		name:  "slice_finder",
		newFn: newSliceFinder,
	},
	{
		name:  "map_finder",
		newFn: newMapFinder,
	},
}

func TestFinderContains(t *testing.T) {
	var (
		numElems = []int{1, 5, 10, 20, 50, 100, 200, 500, 1000}
		hitRates = []float64{1.0, 0.5, 0.25, 0}
	)

	for _, finderInfo := range allFinders {
		t.Run(fmt.Sprintf("finder=%s", finderInfo.name), func(t *testing.T) {
			for _, numElems := range numElems {
				t.Run(fmt.Sprintf("num_elems=%d", numElems), func(t *testing.T) {
					for _, hitRate := range hitRates {
						var (
							f               = finderInfo.newFn()
							addSet, testSet = generateElems(numElems, hitRate)
						)
						for _, elem := range addSet {
							f.add(elem)
						}

						t.Run(fmt.Sprintf("hit_rate=%.2f", hitRate), func(t *testing.T) {
							testFinderContains(t, f, testSet, hitRate)
						})
					}
				})
			}
		})
	}
}

func testFinderContains(t *testing.T, f finder, elems []string, expectedHitRate float64) {
	var (
		hits = 0
	)

	for _, elem := range elems {
		if f.contains(elem) {
			hits++
		}
	}

	hitRate := float64(hits) / float64(len(elems))
	if hitRate != expectedHitRate {
		t.Errorf("unexpected hit rate (expected = %.2f, actual = %.2f)", expectedHitRate, hitRate)
	}
}

var res bool

func BenchmarkFinderContains(b *testing.B) {
	type sets struct {
		addSet  []string
		testSet []string
	}

	var (
		numElems      = []int{1, 5, 10, 20, 50, 100, 200, 500, 1000}
		hitRates      = []float64{1.0, 0.5, 0.25, 0}
		generatedSets = map[int]map[float64]sets{} // used to ensure each implementation gets the same set
	)

	// helper to make sure sets remain the same
	getSets := func(numElems int, hitRate float64) ([]string, []string) {
		if numElemSets, ok := generatedSets[numElems]; ok {
			if hitRateSets, ok := numElemSets[hitRate]; ok {
				return hitRateSets.addSet, hitRateSets.testSet
			}
		} else {
			generatedSets[numElems] = map[float64]sets{}
		}

		addSet, testSet := generateElems(numElems, hitRate)
		generatedSets[numElems][hitRate] = sets{addSet: addSet, testSet: testSet}
		return addSet, testSet
	}

	for _, finderInfo := range allFinders {
		b.Run(fmt.Sprintf("finder=%s", finderInfo.name), func(b *testing.B) {
			for _, numElems := range numElems {
				b.Run(fmt.Sprintf("num_elems=%d", numElems), func(b *testing.B) {
					for _, hitRate := range hitRates {
						var (
							addSet, testSet = getSets(numElems, hitRate)
						)

						f := finderInfo.newFn()

						for _, elem := range addSet {
							f.add(elem)
						}

						b.Run(fmt.Sprintf("hit_rate=%.2f", hitRate), func(b *testing.B) {
							b.ResetTimer()
							benchmarkFinderContains(b, f, testSet)
						})
					}
				})
			}
		})
	}
}

func benchmarkFinderContains(b *testing.B, f finder, elems []string) {
	b.Helper()
	found := false
	for n := 0; n < b.N; n++ {
		for _, elem := range elems {
			found = f.contains(elem)
		}
	}
	res = found
}

func generateElems(numElems int, hitRate float64) ([]string, []string) {
	var (
		addSet     = make([]string, numElems)
		testSetLen int
		r          = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	if hitRate == 0 || hitRate == 1.0 {
		testSetLen = numElems
	} else {
		testSetLen = numElems * int(1/hitRate)
	}
	testSet := make([]string, testSetLen)

	for i := 0; i < numElems; i++ {
		addSet[i] = strconv.Itoa(i)
	}

	switch hitRate {
	case 1:
		copy(testSet, addSet)
	case 0:
		for i := 0; i < numElems; i++ {
			testSet[i] = strconv.Itoa(i + numElems)
		}
	default:
		copy(testSet, addSet)
		for i := numElems; i < testSetLen; i++ {
			testSet[i] = strconv.Itoa(i)
		}
	}

	r.Shuffle(len(addSet), func(i, j int) {
		addSet[i], addSet[j] = addSet[j], addSet[i]
	})
	r.Shuffle(len(testSet), func(i, j int) {
		testSet[i], testSet[j] = testSet[j], testSet[i]
	})

	return addSet, testSet
}
