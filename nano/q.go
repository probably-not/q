package nano

import (
	"sync/atomic"

	"github.com/probably-not/q/consts"
)

type Q uint32

func NewQ() Q {
	return 0
}

// Pop will calculate the position that can currently be popped from the queue.
// It returns the position, a save point (to allow ensuring the commit),
// along with a boolean indicating if the queue is empty or not.
// If the queue is empty, `-1, 0, true` will be returned.
// If the queue is not empty, `pos, savepoint, false` will be returned.
// After receiving the position and savepoint, PopCommit must be called in
// order to ensure that the job is truly the caller's job, and it has not been
// committed by another consumer.
func (q *Q) Pop(factor int) (int, uint32, bool) {
	acquired := atomic.LoadUint32((*uint32)(q))
	mask := (uint32(1) << factor) - 1
	head := acquired & mask
	tail := acquired >> 16 & mask

	if head == tail {
		return -1, 0, true
	}

	return int(tail), acquired, false
}

// PopCommit will commit the previously executed Pop operation to the queue.
// This moves the index of the queue to the next pop-able index.
// It requires a savepoint that was returned by the Pop operation, which will
// be used to ensure that the operation is in fact atomic.
// If the commit succeeds, `true` is returned to indicate that the caller does in fact
// have the job that it received in the Pop operation.
func (q *Q) PopCommit(savepoint uint32) bool {
	return atomic.CompareAndSwapUint32((*uint32)(q), savepoint, uint32(savepoint+consts.CommitPopU32))
}

// Push will calculate the position that can currently be pushed to in the queue.
// It returns the position, along with a boolean indicating if the queue is full or not.
// If the queue is full, `-1, true` will be returned.
// If the queue is not full, `pos, false` will be returned.
func (q *Q) Push(factor int) (int, bool) {
	acquired := atomic.LoadUint32((*uint32)(q))
	mask := (uint32(1) << factor) - 1
	head := acquired & mask
	tail := acquired >> 16 & mask
	next := (head + uint32(1)) & mask

	if acquired&consts.PushOverflowCheckU32 != 0 {
		atomic.AddUint32((*uint32)(q), consts.PushOverflowProtectionU32)
	}

	if next == tail {
		return -1, true
	}

	return int(head), false
}

// PushCommit will commit the previously executed Push operation to the queue.
// This moves the index of the queue to the next push-able index.
func (q *Q) PushCommit() {
	atomic.AddUint32((*uint32)(q), 1)
}
