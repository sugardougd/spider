package terminal

import (
	"container/list"
)

type History struct {
	limit   int
	history *list.List
	current *list.Element
}

func NewHistory(limit int) *History {
	h := History{
		limit:   limit,
		history: list.New(),
	}
	return &h
}

func (h *History) Push(s string) {
	if len(s) == 0 {
		return
	}
	elem := h.history.PushBack(s)
	h.current = elem
	for h.history.Len() > h.limit && h.history.Len() > 0 {
		h.history.Remove(h.history.Front())
	}
}

func (h *History) Prev() (string, bool) {
	if h.current == nil {
		return "", false
	}
	current := h.current.Prev()
	if current == nil {
		return "", false
	}
	h.current = current
	return current.Value.(string), true
}

func (h *History) Next() (string, bool) {
	if h.current == nil {
		return "", false
	}
	current := h.current.Next()
	if current == nil {
		return "", false
	}

	h.current = current
	return current.Value.(string), true
}

func (h *History) Reset() {
	h.history = list.New()
	h.current = nil
}
