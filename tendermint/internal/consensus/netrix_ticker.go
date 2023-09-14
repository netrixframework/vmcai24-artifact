package consensus

import (
	"context"
	"sync"
	"time"

	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/service"

	netrix "github.com/netrixframework/go-clientlibrary"
)

type tiWrapper struct {
	ti timeoutInfo
}

func (ti tiWrapper) Key() string {
	return ti.ti.String()
}

func (ti tiWrapper) Duration() time.Duration {
	return ti.ti.Duration
}

type NetrixTicker struct {
	service.BaseService
	logger      log.Logger
	client      *netrix.ReplicaClient
	timeoutChan chan netrix.TimeoutInfo

	outChan       chan timeoutInfo
	activeTimeout timeoutInfo
	lock          *sync.Mutex
	doneChan      chan struct{}
	once          *sync.Once
}

var _ TimeoutTicker = &NetrixTicker{}
var _ netrix.TimeoutInfo = tiWrapper{}

func NewNetrixTicker(logger log.Logger) (*NetrixTicker, error) {
	client, err := netrix.GetClient()
	if err != nil {
		return nil, err
	}
	t := &NetrixTicker{
		client:      client,
		logger:      logger,
		timeoutChan: client.TimeoutChan(),
		outChan:     make(chan timeoutInfo, 2),
		lock:        new(sync.Mutex),
		doneChan:    make(chan struct{}),
		once:        new(sync.Once),
	}
	t.BaseService = *service.NewBaseService(logger, "NetrixTimeoutTicker", t)
	return t, nil
}

func (n *NetrixTicker) Chan() <-chan timeoutInfo {
	return n.outChan
}

func (n *NetrixTicker) ScheduleTimeout(ti timeoutInfo) {
	n.lock.Lock()
	curTi := n.activeTimeout
	n.lock.Unlock()
	if ti.Height < curTi.Height {
		return
	} else if ti.Height == curTi.Height {
		if ti.Round < curTi.Round {
			return
		} else if ti.Round == curTi.Round {
			if ti.Step > 0 && ti.Step <= curTi.Step {
				return
			}
		}
	}
	n.lock.Lock()
	n.activeTimeout = ti
	n.lock.Unlock()
	if ti.Duration < 0 {
		n.fireTimeout(ti)
	} else {
		n.client.StartTimer(tiWrapper{ti})
	}
}

func (n *NetrixTicker) OnStart(ctx context.Context) error {
	go n.poll(ctx)
	return nil
}

func (n *NetrixTicker) OnStop() {
	n.once.Do(func() {
		close(n.doneChan)
	})
}

func (n *NetrixTicker) fireTimeout(ti timeoutInfo) {
	select {
	case n.outChan <- ti:
	case <-n.doneChan:
	}
}

func (n *NetrixTicker) poll(ctx context.Context) {
	for {
		select {
		case t := <-n.timeoutChan:
			ti := t.(tiWrapper).ti
			n.lock.Lock()
			curTi := n.activeTimeout
			n.lock.Unlock()
			if ti.Height < curTi.Height {
				continue
			} else if ti.Height == curTi.Height {
				if ti.Round < curTi.Round {
					continue
				} else if ti.Round == curTi.Round {
					if ti.Step > 0 && ti.Step < curTi.Step {
						continue
					}
				}
			}
			go n.fireTimeout(ti)
		case <-ctx.Done():
			n.once.Do(func() {
				close(n.doneChan)
			})
			return
		case <-n.doneChan:
			return
		}
	}
}
