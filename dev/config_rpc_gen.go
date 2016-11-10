// DO NOT EDIT. Generated by 'gorums' plugin for protoc-gen-go
// Source file to edit is: config_rpc_tmpl

package dev

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"
)

// ReadReply encapsulates the reply from a Read RPC invocation.
// It contains the id of each node in the quorum that replied and a single
// reply.
type ReadReply struct {
	NodeIDs []uint32
	*State
}

func (r ReadReply) String() string {
	return fmt.Sprintf("node ids: %v | answer: %v", r.NodeIDs, r.State)
}

// Read invokes a Read RPC on configuration c
// and returns the result as a ReadReply.
func (c *Configuration) Read(ctx context.Context, args *ReadRequest) (*ReadReply, error) {
	return c.mgr.read(ctx, c, args)
}

// ReadFuture is a reference to an asynchronous Read RPC invocation.
type ReadFuture struct {
	reply *ReadReply
	err   error
	c     chan struct{}
}

// ReadFuture asynchronously invokes a Read RPC on configuration c and
// returns a ReadFuture which can be used to inspect the RPC reply and error
// when available.
func (c *Configuration) ReadFuture(ctx context.Context, args *ReadRequest) *ReadFuture {
	f := new(ReadFuture)
	f.c = make(chan struct{}, 1)
	go func() {
		defer close(f.c)
		f.reply, f.err = c.mgr.read(ctx, c, args)
	}()
	return f
}

// Get returns the reply and any error associated with the ReadFuture.
// The method blocks until a reply or error is available.
func (f *ReadFuture) Get() (*ReadReply, error) {
	<-f.c
	return f.reply, f.err
}

// Done reports if a reply or error is available for the ReadFuture.
func (f *ReadFuture) Done() bool {
	select {
	case <-f.c:
		return true
	default:
		return false
	}
}

// ReadCorrectable asynchronously invokes a
// correctable Read quorum call on configuration c and returns a
// ReadCorrectable which can be used to inspect any repies or errors
// when available.
func (c *Configuration) ReadCorrectable(ctx context.Context, args *ReadRequest) *ReadCorrectable {
	corr := &ReadCorrectable{
		level:  LevelNotSet,
		donech: make(chan struct{}),
	}
	go func() {
		c.mgr.readCorrectable(ctx, c, corr, args)
	}()
	return corr
}

// ReadCorrectable is a reference to a correctable Read quorum call.
type ReadCorrectable struct {
	mu       sync.Mutex
	reply    *ReadReply
	level    int
	err      error
	done     bool
	watchers []struct {
		level int
		ch    chan struct{}
	}
	donech chan struct{}
}

// Get returns the reply, level and any error associated with the
// ReadCorrectable. The method does not block until a (possibly
// itermidiate) reply or error is available. Level is set to LevelNotSet if no
// reply has yet been received. The Done or Watch methods should be used to
// ensure a reply is available.
func (c *ReadCorrectable) Get() (*ReadReply, int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.reply, c.level, c.err
}

// Done returns a channel that's closed when the Read correctable
// quorum call is done. A call is considered done when the quorum function has
// signaled that a quorum of replies was received or that the call returned an
// error.
func (c *ReadCorrectable) Done() <-chan struct{} {
	return c.donech
}

// Watch returns a channel that's closed when a reply or error at or above the
// specified level is available. If the call is done, the channel is closed
// disregardless of the specified level.
func (c *ReadCorrectable) Watch(level int) <-chan struct{} {
	ch := make(chan struct{})
	c.mu.Lock()
	if level < c.level {
		close(ch)
		c.mu.Unlock()
		return ch
	}
	c.watchers = append(c.watchers, struct {
		level int
		ch    chan struct{}
	}{level, ch})
	c.mu.Unlock()
	return ch
}

func (c *ReadCorrectable) set(reply *ReadReply, level int, err error, done bool) {
	c.mu.Lock()
	if c.done {
		c.mu.Unlock()
		panic("set(...) called on a done correctable")
	}
	c.reply, c.level, c.err, c.done = reply, level, err, done
	if done {
		close(c.donech)
		for _, watcher := range c.watchers {
			close(watcher.ch)
		}
		c.mu.Unlock()
		return
	}
	for i := range c.watchers {
		if c.watchers[i].level <= level {
			close(c.watchers[i].ch)
		}
		c.watchers = c.watchers[i+1:]
	}
	c.mu.Unlock()
}

// WriteReply encapsulates the reply from a Write RPC invocation.
// It contains the id of each node in the quorum that replied and a single
// reply.
type WriteReply struct {
	NodeIDs []uint32
	*WriteResponse
}

func (r WriteReply) String() string {
	return fmt.Sprintf("node ids: %v | answer: %v", r.NodeIDs, r.WriteResponse)
}

// Write invokes a Write RPC on configuration c
// and returns the result as a WriteReply.
func (c *Configuration) Write(ctx context.Context, args *State) (*WriteReply, error) {
	return c.mgr.write(ctx, c, args)
}

// WriteFuture is a reference to an asynchronous Write RPC invocation.
type WriteFuture struct {
	reply *WriteReply
	err   error
	c     chan struct{}
}

// WriteFuture asynchronously invokes a Write RPC on configuration c and
// returns a WriteFuture which can be used to inspect the RPC reply and error
// when available.
func (c *Configuration) WriteFuture(ctx context.Context, args *State) *WriteFuture {
	f := new(WriteFuture)
	f.c = make(chan struct{}, 1)
	go func() {
		defer close(f.c)
		f.reply, f.err = c.mgr.write(ctx, c, args)
	}()
	return f
}

// Get returns the reply and any error associated with the WriteFuture.
// The method blocks until a reply or error is available.
func (f *WriteFuture) Get() (*WriteReply, error) {
	<-f.c
	return f.reply, f.err
}

// Done reports if a reply or error is available for the WriteFuture.
func (f *WriteFuture) Done() bool {
	select {
	case <-f.c:
		return true
	default:
		return false
	}
}

// WriteAsync invokes an asynchronous WriteAsync RPC on configuration c.
// The call has no return value and is invoked on every node in the
// configuration.
func (c *Configuration) WriteAsync(ctx context.Context, args *State) error {
	return c.mgr.writeAsync(ctx, c, args)
}
