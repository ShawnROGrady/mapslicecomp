package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

var allFinders = []struct {
	name  string
	newFn func([]string) finder
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
		numElems = []int{1, 2, 4, 5, 10, 20, 50, 100, 200, 500, 1000}
		hitRates = []float64{1.0, 0.5, 0.25, 0}
		maxLen   = 1000
	)

	if testing.Short() {
		maxLen = 100
	}

	for _, finderInfo := range allFinders {
		t.Run(fmt.Sprintf("finder=%s", finderInfo.name), func(t *testing.T) {
			for _, numElems := range numElems {
				t.Run(fmt.Sprintf("num_elems=%d", numElems), func(t *testing.T) {
					if testing.Short() && numElems > maxLen {
						t.Skip("skipping long test")
					}
					for _, hitRate := range hitRates {
						t.Run(fmt.Sprintf("hit_rate=%.2f", hitRate), func(t *testing.T) {
							numDifferent := float64(numElems) * (1 - hitRate)
							if math.Mod(numDifferent, 1) != 0 {
								t.Skip("skipping due to non-whole num_elems*hit_rate")
							}

							var (
								f               = finderInfo.newFn([]string{})
								addSet, testSet = generateElems(numElems, maxLen, 1, hitRate)
							)
							for _, elem := range addSet {
								f.add(elem)
							}

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

func TestDedupe(t *testing.T) {
	var (
		numElems       = []int{1, 5, 10, 20, 50, 100, 200, 500, 1000}
		duplicateRates = []float64{0.8, 0.6, 0.4, 0.2, 0.0}
	)

	for _, finderInfo := range allFinders {
		t.Run(fmt.Sprintf("finder=%s", finderInfo.name), func(t *testing.T) {
			for _, numElems := range numElems {
				t.Run(fmt.Sprintf("num_elems=%d", numElems), func(t *testing.T) {
					for _, dupRate := range duplicateRates {
						t.Run(fmt.Sprintf("dup_rate=%.2f", dupRate), func(t *testing.T) {
							testDedupe(t, finderInfo.newFn, numElems, dupRate)
						})
					}
				})
			}
		})
	}
}

func testDedupe(t *testing.T, newFinderFn func([]string) finder, numElems int, dupRate float64) {
	var (
		f              = newFinderFn([]string{})
		numDups        = int(float64(numElems) * dupRate)
		expectedUnique = numElems - numDups
	)

	inElems, err := generateElemsWithDups(numElems, dupRate)
	if err != nil {
		t.Fatalf("unexpected error generating test elements: %s", err)
	}

	deduped := dedupe(inElems, f)
	if len(deduped) != expectedUnique {
		t.Errorf("unexpected number of de-deduped elements (expected = %d, actual = %d)", expectedUnique, len(deduped))
	}

	// sort to verify uniqueness
	sort.Strings(deduped)

	for i, elem := range deduped {
		if i == 0 {
			continue
		}

		if elem == deduped[i-1] {
			t.Errorf("unexpected duplicate element found: %s", elem)
		}
	}
}

var res bool

func BenchmarkFinderContains(b *testing.B) {
	var (
		numElems = []int{1, 2, 4, 5, 10, 20, 50, 100, 200, 500, 1000}
		hitRates = []float64{1.0, 0.5, 0.25, 0.2, 0}
		strLens  = []int{1, 2, 5, 10, 20, 30, 50, 100}
		maxLen   = 1000
	)

	if testing.Short() {
		maxLen = 100
	}

	for _, finderInfo := range allFinders {
		b.Run(fmt.Sprintf("finder=%s", finderInfo.name), func(b *testing.B) {
			for _, numElems := range numElems {
				b.Run(fmt.Sprintf("num_elems=%d", numElems), func(b *testing.B) {
					if testing.Short() && numElems > maxLen {
						b.Skip("skipping long benchmark")
					}
					for _, hitRate := range hitRates {
						b.Run(fmt.Sprintf("hit_rate=%.2f", hitRate), func(b *testing.B) {
							for _, strLen := range strLens {
								b.Run(fmt.Sprintf("str_len=%d", strLen), func(b *testing.B) {
									numDifferent := float64(numElems) * (1 - hitRate)
									if math.Mod(numDifferent, 1) != 0 {
										b.Skip("skipping due to non-whole num_elems*hit_rate")
									}
									benchmarkFinderContains(b, finderInfo.newFn, numElems, maxLen, strLen, hitRate)
								})
							}
						})
					}
				})
			}
		})
	}
}

func benchmarkFinderContains(b *testing.B, newFinderFn func([]string) finder, numElems, maxLen, strLen int, hitRate float64) {
	b.Helper()
	found := false
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		addSet, testSet := generateElems(numElems, maxLen, strLen, hitRate)
		f := newFinderFn(addSet)
		b.StartTimer()

		for _, elem := range testSet {
			found = f.contains(elem)
		}
	}
	res = found
}

func generateElems(numElems, testSetLen, strLen int, hitRate float64) ([]string, []string) {
	var (
		addSet       = make([]string, numElems)
		testElems    = make([]string, numElems)
		testSet      = make([]string, testSetLen)
		r            = rand.New(rand.NewSource(time.Now().UnixNano()))
		numDifferent = int(float64(numElems) * (1 - hitRate))
	)

	for i := 0; i < numElems; i++ {
		addSet[i] = strings.Repeat(strconv.Itoa(i), strLen)
	}

	switch hitRate {
	case 1:
		copy(testElems, addSet)
	case 0:
		for i := 0; i < numElems; i++ {
			testElems[i] = strings.Repeat(strconv.Itoa(i+numElems), strLen)
		}
	default:
		copy(testElems, addSet)
		toChange := r.Perm(numElems)[:numDifferent]
		for i, diffI := range toChange {
			testElems[diffI] = strings.Repeat(strconv.Itoa(i+numElems), strLen)
		}
	}

	// bring testSet up to proper len
	for i := 0; i < testSetLen; i += numElems {
		copy(testSet[i:], testElems)
	}

	r.Shuffle(len(addSet), func(i, j int) {
		addSet[i], addSet[j] = addSet[j], addSet[i]
	})
	r.Shuffle(len(testSet), func(i, j int) {
		testSet[i], testSet[j] = testSet[j], testSet[i]
	})

	return addSet, testSet
}

var dedupeRes []string

func BenchmarkDedupe(b *testing.B) {
	var (
		numElems       = []int{1, 2, 3, 4, 5, 10, 20, 50, 100, 200, 500, 1000}
		maxLen         = 1000
		duplicateRates = []float64{0.8, 0.6, 0.4, 0.2, 0.0}
	)

	if testing.Short() {
		maxLen = 100
	}

	for _, finderInfo := range allFinders {
		b.Run(fmt.Sprintf("finder=%s", finderInfo.name), func(b *testing.B) {
			for _, numElems := range numElems {
				b.Run(fmt.Sprintf("num_elems=%d", numElems), func(b *testing.B) {
					if testing.Short() && numElems > maxLen {
						b.Skip("skipping long benchmark")
					}
					for _, dupRate := range duplicateRates {
						b.Run(fmt.Sprintf("dup_rate=%.2f", dupRate), func(b *testing.B) {
							var (
								numDups = float64(numElems) * dupRate
							)
							if math.Mod(numDups, 1) != 0 {
								b.Skip("skipping due to non-whole num_elems*dup_rate")
							}
							benchmarkDedupe(b, finderInfo.newFn, numElems, dupRate)
						})
					}
				})
			}
		})
	}
}

func benchmarkDedupe(b *testing.B, newFinderFn func([]string) finder, numElems int, dupRate float64) {
	b.Helper()
	deduped := []string{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		f := newFinderFn([]string{})
		elems, err := generateElemsWithDups(numElems, dupRate)
		if err != nil {
			b.Fatalf("unexpected error generating test elems: %s", err)
		}
		b.StartTimer()

		deduped = dedupe(elems, f)
	}
	dedupeRes = deduped
}

func generateElemsWithDups(numElems int, dupRate float64) ([]string, error) {
	var (
		numDups   = int(float64(numElems) * dupRate)
		numUnique = numElems - numDups
		elems     = make([]string, numElems)
		r         = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	if numUnique <= 0 {
		return elems, fmt.Errorf("numUnique = %d with numElems = %d, dupRate = %.2f", numUnique, numElems, dupRate)
	}

	// unique elements
	for i := 0; i < numUnique; i++ {
		elems[i] = strconv.Itoa(i)
	}

	// duplicate elements
	for i := numUnique; i < numElems; i++ {
		elems[i] = strconv.Itoa(r.Intn(numUnique))
	}

	r.Shuffle(len(elems), func(i, j int) {
		elems[i], elems[j] = elems[j], elems[i]
	})

	return elems, nil
}
