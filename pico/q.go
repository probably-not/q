package pico

import "sync/atomic"

type Q uint32

const (
	release               = uint32(0x10000)
	overflowCheck         = uint32(0x8000)
	overflowProtectionBit = uint32(-0x8000 & 0xffffffff)
)

// Pop will calculate the position that can currently be popped from the queue.
// It returns the position, along with a boolean indicating if the queue is empty or not.
// If the queue is empty, `-1, true` will be returned.
// If the queue is not empty, `pos, false` will be returned.
func Pop(q *Q, exp int) (int, bool) {
	acquired := atomic.LoadUint32((*uint32)(q))
	mask := (uint32(1) << exp) - 1
	head := acquired & mask
	tail := acquired >> 16 & mask

	if head == tail {
		return -1, true
	}

	return int(tail), false
}

// PopCommit will commit the previously executed Pop operation to the queue.
// This moves the index of the queue to the next pop-able index.
func PopCommit(q *Q) {
	atomic.AddUint32((*uint32)(q), release)
}

// Push will calculate the position that can currently be pushed to in the queue.
// It returns the position, along with a boolean indicating if the queue is full or not.
// If the queue is full, `-1, true` will be returned.
// If the queue is not full, `pos, false` will be returned.
func Push(q *Q, exp int) (int, bool) {
	acquired := atomic.LoadUint32((*uint32)(q))
	mask := (uint32(1) << exp) - 1
	head := acquired & mask
	tail := acquired >> 16 & mask
	next := (head + uint32(1)) & mask

	if acquired&overflowCheck != 0 {
		atomic.AddUint32((*uint32)(q), overflowProtectionBit)
	}

	if next == tail {
		return -1, true
	}

	return int(tail), false
}

// PushCommit will commit the previously executed Push operation to the queue.
// This moves the index of the queue to the next push-able index.
func PushCommit(queue *Q) {
	atomic.AddUint32((*uint32)(queue), 1)
}
