package spider

import "fmt"

// RingBufferHistory is a ring buffer of strings.
type RingBufferHistory struct {
	// Max number of elements,must greater than zero
	Max int
	// entries contains Max elements.
	entries []string
	// head contains the index of the element most recently added to the ring.
	head int
	// size contains the number of elements in the ring.
	size int
}

func (s *RingBufferHistory) Add(a string) {
	if s.entries == nil {
		if s.Max < 1 {
			s.Max = 100
		}
		s.entries = make([]string, s.Max)
	}
	if len(a) == 0 {
		return
	}

	s.head = (s.head + 1) % s.Max
	s.entries[s.head] = a
	if s.size < s.Max {
		s.size++
	}
}

func (s *RingBufferHistory) Len() int {
	return s.size
}

// At returns the value passed to the nth previous call to Add.
// If n is zero then the immediately prior value is returned, if one, then the
// next most recent, and so on. If such an element doesn't exist then ok is
// false.
func (s *RingBufferHistory) At(n int) string {
	if n < 0 || n >= s.size {
		panic(fmt.Sprintf("term: history index [%d] out of range [0,%d)", n, s.size))
	}
	index := s.head - n
	if index < 0 {
		index += s.Max
	}
	return s.entries[index]
}
