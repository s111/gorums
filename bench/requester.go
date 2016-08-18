package bench

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	rpc "github.com/relab/gorums/dev"
	"github.com/tylertreat/bench"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func init() {
	silentLogger := log.New(ioutil.Discard, "", log.LstdFlags)
	grpclog.SetLogger(silentLogger)
	grpc.EnableTracing = false
	rand.Seed(42)
}

type RequesterFactory struct {
	Addrs             []string
	ReadQuorum        int
	PayloadSize       int
	QRPCTimeout       time.Duration
	WriteRatioPercent int
}

func (r *RequesterFactory) GetRequester(uint64) bench.Requester {
	return &gorumsRequester{
		addrs: r.Addrs,
		dialOpts: []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(time.Second),
		},
		readq:       r.ReadQuorum,
		payloadSize: r.PayloadSize,
		timeout:     r.QRPCTimeout,
		writeRatio:  r.WriteRatioPercent,
	}
}

type gorumsRequester struct {
	addrs       []string
	dialOpts    []grpc.DialOption
	readq       int
	payloadSize int
	timeout     time.Duration
	writeRatio  int

	mgr    *rpc.Manager
	config *rpc.Configuration
	state  *rpc.State
}

func (gr *gorumsRequester) Setup() error {
	var err error
	gr.mgr, err = rpc.NewManager(
		gr.addrs,
		rpc.WithGrpcDialOptions(gr.dialOpts...),
	)
	if err != nil {
		return err
	}

	ids := gr.mgr.NodeIDs()
	qspec := newRegisterQSpec(gr.readq, len(gr.addrs))
	gr.config, err = gr.mgr.NewConfiguration(ids, qspec, gr.timeout)
	if err != nil {
		return err
	}

	// Set initial state.
	gr.state = &rpc.State{
		Value:     strings.Repeat("x", gr.payloadSize),
		Timestamp: time.Now().UnixNano(),
	}
	wreply, err := gr.config.Write(gr.state)
	if err != nil {
		return fmt.Errorf("write rpc error: %v", err)
	}
	if !wreply.Reply.New {
		return fmt.Errorf("intital write reply was not marked as new")
	}

	return nil
}

func (gr *gorumsRequester) Request() error {
	var err error
	switch gr.writeRatio {
	case 0:
		_, err = gr.config.Read(&rpc.ReadRequest{})
	case 100:
		gr.state.Timestamp = time.Now().UnixNano()
		_, err = gr.config.Write(gr.state)
	default:
		x := rand.Intn(100)
		if x < gr.writeRatio {
			gr.state.Timestamp = time.Now().UnixNano()
			_, err = gr.config.Write(gr.state)
		} else {
			_, err = gr.config.Read(&rpc.ReadRequest{})
		}
	}

	return err
}

func (gr *gorumsRequester) Teardown() error {
	gr.mgr.Close()
	gr.mgr = nil
	gr.config = nil
	return nil
}

type registerQSpec struct {
	rq, wq int
}

func newRegisterQSpec(rq, wq int) rpc.QuorumSpec {
	return &registerQSpec{
		rq: rq,
		wq: wq,
	}
}

func (rqs *registerQSpec) ReadQF(replies []*rpc.State) (*rpc.State, bool) {
	if len(replies) < rqs.rq {
		return nil, false
	}
	return replies[0], true
}

func (rqs *registerQSpec) WriteQF(replies []*rpc.WriteResponse) (*rpc.WriteResponse, bool) {
	if len(replies) < rqs.wq {
		return nil, false
	}
	return replies[0], true
}