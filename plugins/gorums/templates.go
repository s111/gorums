// DO NOT EDIT. Generated by github.com/relab/gorums/cmd/gentemplates
// Template source files to edit is in the: 'dev' folder

package gorums

const config_rpc_tmpl = `
{{/* Remember to run 'make gengolden' after editing this file. */}}

{{- if not .IgnoreImports}}
package {{.PackageName}}

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"
)

{{- end}}

{{range $elm := .Services}}

{{if .Multicast}}

// {{.MethodName}} is a one-way multicast operation, where args is sent to
// every node in configuration c. The call is asynchronous and has no response
// return value.
func (c *Configuration) {{.MethodName}}(ctx context.Context, args *{{.ReqName}}) error {
	return c.mgr.{{.UnexportedMethodName}}(ctx, c, args)
}

{{- end -}}

{{if .QuorumCall}}

// {{.TypeName}} encapsulates the reply from a {{.MethodName}} quorum call.
// It contains the id of each node of the quorum that replied and a single reply.
type {{.TypeName}} struct {
	NodeIDs []uint32
	*{{.RespName}}
}

func (r {{.TypeName}}) String() string {
	return fmt.Sprintf("node ids: %v | answer: %v", r.NodeIDs, r.{{.RespName}})
}

// {{.MethodName}} invokes a {{.MethodName}} quorum call on configuration c
// and returns the result as a {{.TypeName}}.
func (c *Configuration) {{.MethodName}}(ctx context.Context, args *{{.ReqName}}) (*{{.TypeName}}, error) {
	return c.mgr.{{.UnexportedMethodName}}(ctx, c, args)
}
{{- end -}}

{{if .Future}}

// {{.MethodName}}Future is a reference to an asynchronous {{.MethodName}} quorum call invocation.
type {{.MethodName}}Future struct {
	reply *{{.TypeName}}
	err   error
	c     chan struct{}
}

// {{.MethodName}}Future asynchronously invokes a {{.MethodName}} quorum call
// on configuration c and returns a {{.MethodName}}Future which can be used to
// inspect the quorum call reply and error when available.
func (c *Configuration) {{.MethodName}}Future(ctx context.Context, args *{{.ReqName}}) *{{.MethodName}}Future {
	f := new({{.MethodName}}Future)
	f.c = make(chan struct{}, 1)
	go func() {
		defer close(f.c)
		f.reply, f.err = c.mgr.{{.UnexportedMethodName}}(ctx, c, args)
	}()
	return f
}

// Get returns the reply and any error associated with the {{.MethodName}}Future.
// The method blocks until a reply or error is available.
func (f *{{.MethodName}}Future) Get() (*{{.TypeName}}, error) {
	<-f.c
	return f.reply, f.err
}

// Done reports if a reply and/or error is available for the {{.MethodName}}Future.
func (f *{{.MethodName}}Future) Done() bool {
	select {
	case <-f.c:
		return true
	default:
		return false
	}
}

{{- end -}}

{{if .Correctable}}

// {{.MethodName}}Correctable asynchronously invokes a
// correctable {{.MethodName}} quorum call on configuration c and returns a
// {{.MethodName}}Correctable which can be used to inspect any repies or errors
// when available.
func (c *Configuration) {{.MethodName}}Correctable(ctx context.Context, args *ReadRequest) *{{.MethodName}}Correctable {
	corr := &{{.MethodName}}Correctable{
		level:  LevelNotSet,
		donech: make(chan struct{}),
	}
	go func() {
		c.mgr.{{.UnexportedMethodName}}Correctable(ctx, c, corr, args)
	}()
	return corr
}

// {{.MethodName}}Correctable is a reference to a correctable {{.MethodName}} quorum call.
type {{.MethodName}}Correctable struct {
	mu       sync.Mutex
	reply    *{{.TypeName}}
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
// {{.MethodName}}Correctable. The method does not block until a (possibly
// itermidiate) reply or error is available. Level is set to LevelNotSet if no
// reply has yet been received. The Done or Watch methods should be used to
// ensure that a reply is available.
func (c *{{.MethodName}}Correctable) Get() (*{{.TypeName}}, int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.reply, c.level, c.err
}

// Done returns a channel that's closed when the {{.MethodName}} correctable
// quorum call is done. A call is considered done when the quorum function has
// signaled that a quorum of replies was received or that the call returned an
// error.
func (c *{{.MethodName}}Correctable) Done() <-chan struct{} {
	return c.donech
}

// Watch returns a channel that's closed when a reply or error at or above the
// specified level is available. If the call is done, the channel is closed
// disregardless of the specified level.
func (c *{{.MethodName}}Correctable) Watch(level int) <-chan struct{} {
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

func (c *{{.MethodName}}Correctable) set(reply *{{.TypeName}}, level int, err error, done bool) {
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

{{- end -}}
{{- end -}}
`

const mgr_rpc_tmpl = `
{{/* Remember to run 'make gengolden' after editing this file. */}}
{{$pkgName := .PackageName}}

{{if not .IgnoreImports}}
package {{$pkgName}}

import (
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)
{{end}}

{{range $elm := .Services}}

{{if .Multicast}}
func (m *Manager) {{.UnexportedMethodName}}(ctx context.Context, c *Configuration, args *{{.ReqName}}) error {
	for _, node := range c.nodes {
		go func(n *Node) {
			err := n.{{.MethodName}}Client.Send(args)
			if err == nil {
				return
			}
			if m.logger != nil {
				m.logger.Printf("%d: {{.UnexportedMethodName}} stream send error: %v", n.id, err)
			}
		}(node)
	}

	return nil
}

{{- end -}}
{{if .QuorumCall}}

type {{.UnexportedTypeName}} struct {
	nid   uint32
	reply *{{.RespName}}
	err   error
}

func (m *Manager) {{.UnexportedMethodName}}(ctx context.Context, c *Configuration, args *{{.ReqName}}) (*{{.TypeName}}, error) {
	replyChan := make(chan {{.UnexportedTypeName}}, c.n)
	newCtx, cancel := context.WithCancel(ctx)

	for _, n := range c.nodes {
		go callGRPC{{.MethodName}}(newCtx, n, args, replyChan)
	}

	var (
		replyValues = make([]*{{.RespName}}, 0, c.n)
		reply       = &{{.TypeName}}{NodeIDs: make([]uint32, 0, c.n)}
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errCount++
				break
			}
			replyValues = append(replyValues, r.reply)
			reply.NodeIDs = append(reply.NodeIDs, r.nid)
			if reply.{{.RespName}}, quorum = c.qspec.{{.MethodName}}QF(replyValues); quorum {
				cancel()
				return reply, nil
			}
		case <-newCtx.Done():
			return reply, QuorumCallError{ctx.Err().Error(), errCount, len(replyValues)}
		}

		if errCount+len(replyValues) == c.n {
			cancel()
			return reply, QuorumCallError{"incomplete call", errCount, len(replyValues)}
		}
	}
}

func callGRPC{{.MethodName}}(ctx context.Context, node *Node, args *{{.ReqName}}, replyChan chan<- {{.UnexportedTypeName}}) {
	reply := new({{.RespName}})
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/{{$pkgName}}.{{.ServName}}/{{.MethodName}}",
		args,
		reply,
		node.conn,
	)
	switch grpc.Code(err) { // nil -> codes.OK
	case codes.OK, codes.Canceled:
		node.setLatency(time.Since(start))
	default:
		node.setLastErr(err)
	}
	replyChan <- {{.UnexportedTypeName}}{node.id, reply, err}
}

{{- end -}}
{{if .Correctable}}

func (m *Manager) {{.UnexportedMethodName}}Correctable(ctx context.Context, c *Configuration, corr *{{.MethodName}}Correctable, args *{{.ReqName}}) {
	replyChan := make(chan {{.UnexportedTypeName}}, c.n)
	newCtx, cancel := context.WithCancel(ctx)

	for _, n := range c.nodes {
		go callGRPC{{.MethodName}}(newCtx, n, args, replyChan)
	}

	var (
		replyValues     = make([]*{{.RespName}}, 0, c.n)
		reply           = &{{.TypeName}}{NodeIDs: make([]uint32, 0, c.n)}
		clevel      	= LevelNotSet
		rlevel      int
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errCount++
				break
			}
			replyValues = append(replyValues, r.reply)
			reply.NodeIDs = append(reply.NodeIDs, r.nid)
			reply.{{.RespName}}, rlevel, quorum = c.qspec.{{.MethodName}}CorrectableQF(replyValues)
			if quorum {
				cancel()
				corr.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				corr.set(reply, rlevel, nil, false)
			}
		case <-newCtx.Done():
			corr.set(reply, clevel, QuorumCallError{ctx.Err().Error(), errCount, len(replyValues)}, true)
			return
		}

		if errCount+len(replyValues) == c.n {
			cancel()
			corr.set(reply, clevel, QuorumCallError{"incomplete call", errCount, len(replyValues)}, true)
			return
		}
	}
}

{{- end -}}
{{- end -}}
`

const node_tmpl = `
{{/* Remember to run 'make gengolden' after editing this file. */}}

{{- if not .IgnoreImports}}
package {{.PackageName}}

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
)
{{- end}}

// Node encapsulates the state of a node on which a remote procedure call
// can be made.
type Node struct {
	// Only assigned at creation.
	id   uint32
	self bool
	addr string
	conn *grpc.ClientConn


{{range .Clients}}
	{{.}} {{.}}
{{end}}

{{range .Services}}
{{if .Multicast}}
	{{.MethodName}}Client {{.ServName}}_{{.MethodName}}Client
{{end}}
{{end}}

	sync.Mutex
	lastErr error
	latency time.Duration
}

func (n *Node) connect(opts ...grpc.DialOption) error {
  	var err error
	n.conn, err = grpc.Dial(n.addr, opts...)
	if err != nil {
		return fmt.Errorf("dialing node failed: %v", err)
	}

{{range .Clients}}
	n.{{.}} = New{{.}}(n.conn)
{{end}}

{{range .Services}}
{{if .Multicast}}
  	n.{{.MethodName}}Client, err = n.{{.ServName}}Client.{{.MethodName}}(context.Background())
  	if err != nil {
  		return fmt.Errorf("stream creation failed: %v", err)
  	}
{{end}}
{{end -}}

	return nil
}

func (n *Node) close() error {
	// TODO: Log error, mainly care about the connection error below.
        // We should log this error, but we currently don't have access to the
        // logger in the manager.
{{- range .Services -}}
{{if .Multicast}}
	_, _ = n.{{.MethodName}}Client.CloseAndRecv()
{{- end -}}
{{end}}
	
	if err := n.conn.Close(); err != nil {
                return fmt.Errorf("conn close error: %v", err)
        }	
	return nil
}
`

const qspec_tmpl = `
{{/* Remember to run 'make gengolden' after editing this file. */}}

{{- if not .IgnoreImports}}
package {{.PackageName}}
{{- end}}

// QuorumSpec is the interface that wraps every quorum function.
type QuorumSpec interface {
{{- range $elm := .Services}}
{{- if .QuorumCall}}
	// {{.MethodName}}QF is the quorum function for the {{.MethodName}}
	// quorum call method.
	{{.MethodName}}QF(replies []*{{.RespName}}) (*{{.RespName}}, bool)
{{end}}

{{if .Correctable}}
	// {{.MethodName}}CorrectableQF is the quorum function for the {{.MethodName}}
	// correctable quorum call method.
	{{.MethodName}}CorrectableQF(replies []*{{.RespName}}) (*{{.RespName}}, int, bool)
{{end}}
{{- end -}}
}
`

var templates = map[string]string{
	"config_rpc_tmpl": config_rpc_tmpl,
	"mgr_rpc_tmpl":    mgr_rpc_tmpl,
	"node_tmpl":       node_tmpl,
	"qspec_tmpl":      qspec_tmpl,
}
