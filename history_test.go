package spider

import "testing"

func TestRingBufferHistory(t *testing.T) {
	history := &RingBufferHistory{}
	t.Run("TestRingBufferHistory", func(t *testing.T) {
		// add hello
		history.Add("hello")
		if len := history.Len(); len != 1 {
			t.Errorf("Len() = %v, want %v", len, 1)
		}
		if s := history.At(0); s != "hello" {
			t.Errorf("At(0) = %s, want %s", s, "hello")
		}

		// add empty string
		history.Add("")
		if len := history.Len(); len != 1 {
			t.Errorf("Len() = %v, want %v", len, 1)
		}
		if s := history.At(0); s != "hello" {
			t.Errorf("At(0) = %s, want %s", s, "hello")
		}

		// add world
		history.Add("world")
		if len := history.Len(); len != 2 {
			t.Errorf("Len() = %v, want %v", len, 2)
		}
		if s := history.At(0); s != "world" {
			t.Errorf("At(0) = %s, want %s", s, "world")
		}
		if s := history.At(1); s != "hello" {
			t.Errorf("At(1) = %s, want %s", s, "hello")
		}
	})
}
