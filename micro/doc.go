// Package micro contains the third most minimalistic version of the implementation of this queue.
// It enables a single producer/multiple consumer architecture to work on the same slice of jobs,
// with the indices of jobs pushed to the queue and jobs popped from the queue fully managed by the `micro.Q`.
// The slice of jobs itself is managed by an outside source, however access to this slice of jobs should be fully
// managed by the `micro.Q` type.
// In the more minimalistic implementations, the size factor of the queue is something that needs to be held outside of
// the queue, and the queue simply receives the size factor from the caller. While this allows for an extremely minimalistic
// implementation, it means that the size factor must be passed to all of the callers of the queue, in order to ensure
// that everyone has the correct factor. This can lead to issues with accidental losses of the factor.
// In this implementation, we sacrifice minimalism and size with usability, and we convert the queue from a raw `uint32`
// to a struct, which will hold the queue itself, along with the size factor.
// Package micro is not safe to use in a multiple producer/multiple consumer scenario, however with a single producer
// you may have multiple consumers.
package micro
