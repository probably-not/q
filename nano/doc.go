// Package nano contains the second most minimalistic version of the implementation of this queue.
// It enables a single producer/multiple consumer architecture to work on the same slice of jobs,
// with the indices of jobs pushed to the queue and jobs popped from the queue fully managed by the `nano.Q`.
// The slice of jobs itself is managed by an outside source, however access to this slice of jobs should be fully
// managed by the `nano.Q` type.
// A caveat of allowing multiple consumers is that the slice of jobs must be wrapped for threadsafe access. This is due to
// the implicit availability of two consumers to read the same job, within the same time that the producer has made a commit.
// In this situation, one of the consumers will commit accurately, while one will read the job, then fail the commit. While this
// is safe for the consumers (since we are not running jobs if the commit fails), this does cause a data race, since the consumer that
// is going to fail will access the same index that the producer is writing concurrently writing to. Because of this limitation, access to
// the slice of jobs must be threadsafe in order to avoid this race condition.
// Package nano is not safe to use in a multiple producer/multiple consumer scenario, however with a single producer
// you may have multiple consumers.
package nano
