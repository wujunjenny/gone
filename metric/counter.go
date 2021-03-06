package metric

import (
	"github.com/One-com/gone/metric/num64"
	"sync/atomic"
)

// Counter is different from a GaugeInt64 in that it is reset to zero every
// time its flushed - and thus being server-side maintained.
type Counter struct {
	name string
	val  int64
}

// NewCounter returns a client side buffered counter (a counter is a server side maintained value).
// "Server side" meaning that it's reset to 0 every time it's sent to the server and the total
// tally is kept on the server.
// This poses the risk of the server-side absolute value to drift in case of increments
// lost in transit. However, it also allows several distributed processes to update the same counter.
// If you want to have a pure client side counter, use GaugeInt64
func NewCounter(name string) *Counter {
	g := &Counter{name: name}
	return g
}

// FlushReading flushes the accumulated counter value to the supplied Sink
func (c *Counter) FlushReading(s Sink) {
	val := atomic.SwapInt64(&c.val, 0)
	if val != 0 {
		n := num64.FromInt64(int64(val))
		s.RecordNumeric64(MeterCounter, c.name, n)
	}
}

// Name returns the name of the counter
func (c *Counter) Name() string {
	return c.name
}

// Inc increased the counter by the supplied value
func (c *Counter) Inc(i int64) {
	atomic.AddInt64(&c.val, i)
}

// Dec decreased the counter by the supplied value
func (c *Counter) Dec(i int64) {
	atomic.AddInt64(&c.val, -i)
}
