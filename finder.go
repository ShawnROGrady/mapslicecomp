package main

import "strings"

type finder interface {
	add(string)
	contains(string) bool
}

type sliceFinder struct {
	sl []string
}

func newSliceFinder(initElems []string) finder {
	var (
		sl = make([]string, len(initElems))
	)
	copy(sl, initElems)
	return &sliceFinder{
		sl: sl,
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

func newMapFinder(initElems []string) finder {
	var (
		m = make(map[string]struct{}, len(initElems))
	)
	for _, elem := range initElems {
		m[elem] = struct{}{}
	}
	return &mapFinder{
		m: m,
	}
}

func (f *mapFinder) add(s string) {
	f.m[s] = struct{}{}
}

func (f *mapFinder) contains(s string) bool {
	_, ok := f.m[s]
	return ok
}

func dedupe(in []string, f finder) []string {
	deduped := []string{}
	for _, s := range in {
		if f.contains(s) {
			continue
		}
		deduped = append(deduped, s)
		f.add(s)
	}
	return deduped
}
