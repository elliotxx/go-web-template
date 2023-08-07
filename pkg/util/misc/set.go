package misc

import (
	"sort"
	"strings"
	"sync"
)

// Set is a helper type for storing unique values in a map.
type Set struct {
	mu    *sync.Mutex
	items map[string]struct{}
}

// NewSet returns an empty sets map.
func NewSet() *Set {
	return &Set{
		mu:    &sync.Mutex{},
		items: map[string]struct{}{},
	}
}

// Set adds a value to the sets maps.
func (st *Set) Set(key string) {
	// st[key] = struct{}{}
	st.mu.Lock()
	defer st.mu.Unlock()
	st.items[key] = struct{}{}
}

// String returns the sets map as a comma-separated string.
func (st *Set) String() string {
	st.mu.Lock()
	defer st.mu.Unlock()

	s := []string{}
	for k := range st.items {
		s = append(s, k)
	}
	sort.Strings(s)

	return strings.Join(s, ",")
}
