package pico

import "testing"

func TestPop(t *testing.T) {
	var queue Q
	queueSize := 6
	idx, isEmpty := Pop(&queue, queueSize)
	if idx != -1 {
		t.Errorf("expected the returned index to be -1, got %d", idx)
	}

	if !isEmpty {
		t.Error("expected isEmpty to be true, got false")
	}
}
