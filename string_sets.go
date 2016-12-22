package aph

import "sync"

type stringSet map[string]bool

func newStringSet(items ...string) stringSet {
	set := stringSet{}

	for _, i := range items {
		set[i] = true
	}

	return set
}

type concurrentStringSet struct {
	set  stringSet
	lock sync.RWMutex
}

func newConcurrentStringSet(items ...string) *concurrentStringSet {
	cset := &concurrentStringSet{
		set: stringSet{},
	}

	for _, i := range items {
		cset.set[i] = true
	}

	return cset
}

func (cset *concurrentStringSet) add(s string) {
	cset.lock.Lock()
	defer cset.lock.Unlock()

	cset.set[s] = true
}

func (cset *concurrentStringSet) has(s string) bool {
	cset.lock.RLock()
	defer cset.lock.RUnlock()

	return cset.set[s]
}

func (cset *concurrentStringSet) remove(s string) {
	cset.lock.Lock()
	defer cset.lock.Unlock()

	delete(cset.set, s)
}

func (cset *concurrentStringSet) slice() []string {
	cset.lock.RLock()
	defer cset.lock.RUnlock()

	slice := []string{}

	for s := range cset.set {
		slice = append(slice, s)
	}

	return slice
}
