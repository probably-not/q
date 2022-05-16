// Package nano contains the second most minimalistic version of the implementation of this queue.
// It enables a single producer/multiple consumer architecture to work on the same slice of jobs,
// with the indices of jobs pushed to the queue and jobs popped from the queue fully managed by the `nano.Q`.
// The slice of jobs itself is managed by an outside source, however access to this slice of jobs should be fully
// managed by the `nano.Q` type.
// Package nano is not safe to use in a multiple producer/multiple consumer scenario, however with a single producer
// you may have multiple consumers.
package nano

const (
	commitPop              = uint32(0x10000)
	pushOverflowCheck      = uint32(0x8000)
	pushOverflowProtection = uint32(-0x8000 & 0xffffffff)
)
