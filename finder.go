package main

import "strings"

type finder interface {
	add(string)
	contains(string) bool
}

type sliceFinder struct {
	sl []string
}

func newSliceFinder() finder {
	return &sliceFinder{
		sl: []string{},
	}
}

func (f *sliceFinder) add(s string) {
	f.sl = append(f.sl, s)
}

func (f *sliceFinder) contains(s string) bool {
	for _, elem := range f.sl {
		if strings.Compare(elem, s) == 0 {
			return true
		}
	}
	return false
}

type mapFinder struct {
	m map[string]struct{}
}

func newMapFinder() finder {
	return &mapFinder{
		m: map[string]struct{}{},
	}
}

func (f *mapFinder) add(s string) {
	f.m[s] = struct{}{}
}

func (f *mapFinder) contains(s string) bool {
	_, ok := f.m[s]
	return ok
}
