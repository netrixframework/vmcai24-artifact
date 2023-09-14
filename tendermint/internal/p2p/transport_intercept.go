package p2p

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	netrix "github.com/netrixframework/go-clientlibrary"
	ntypes "github.com/netrixframework/netrix/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/internal/p2p/conn"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/service"
	"github.com/tendermint/tendermint/types"
)

type router struct {
	connections map[types.NodeID]*InterceptedConnection
	lock        *sync.Mutex
}

func newRouter() *router {
	return &router{
		connections: make(map[types.NodeID]*InterceptedConnection),
		lock:        new(sync.Mutex),
	}
}

func (r *router) addConnection(nodeID types.NodeID, conn *InterceptedConnection) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.connections[nodeID] = conn
}

func (r *router) dispatchMessage(m *ntypes.Message) {
	nID := types.NodeID(m.From)
	r.lock.Lock()
	ch, ok := r.connections[nID]
	r.lock.Unlock()
	if ok {
		select {
		case <-ch.doneCh:
		default:
			ch.recvCh <- m
		}
	}
}

type InterceptedTransportOptions struct {
	MConnOptions MConnTransportOptions
	NetrixConfig *netrix.Config
}

// InterceptedTransport wraps around MConnTransport
type InterceptedTransport struct {
	service.BaseService

	nodeID       string
	netrixClient *netrix.ReplicaClient
	router       *router
	logger       log.Logger

	mConnTransport *MConnTransport
	mConnConfig    conn.MConnConfig
	channelDescs   []*ChannelDescriptor
	mConnOptions   MConnTransportOptions
	mConnLock      *sync.Mutex

	doneCh chan struct{}
	once   *sync.Once
}

var _ Transport = &InterceptedTransport{}

// NewInterceptedTransport create InterceptedTransport
func NewInterceptedTransport(
	logger log.Logger,
	mConnConfig conn.MConnConfig,
	channelDescs []*ChannelDescriptor,
	options InterceptedTransportOptions,
	directiveHandler netrix.DirectiveHandler,
) (*InterceptedTransport, error) {
	err := netrix.Init(options.NetrixConfig, directiveHandler, logger)
	if err != nil {
		return nil, err
	}
	client, _ := netrix.GetClient()
	t := &InterceptedTransport{
		logger: logger,
		mConnTransport: NewMConnTransport(
			logger,
			mConnConfig,
			channelDescs,
			options.MConnOptions,
		),
		nodeID:       string(options.NetrixConfig.ReplicaID),
		mConnConfig:  mConnConfig,
		channelDescs: channelDescs,
		mConnOptions: options.MConnOptions,
		mConnLock:    new(sync.Mutex),
		netrixClient: client,
		router:       newRouter(),
		doneCh:       make(chan struct{}),
		once:         new(sync.Once),
	}
	t.BaseService = *service.NewBaseService(logger, "netrix-transport", t)
	return t, nil
}

func (t *InterceptedTransport) String() string {
	return "netrix"
}

func (t *InterceptedTransport) Protocols() []Protocol {
	t.mConnLock.Lock()
	defer t.mConnLock.Unlock()

	return t.mConnTransport.Protocols()
}

func (t *InterceptedTransport) Endpoint() (*Endpoint, error) {
	t.mConnLock.Lock()
	defer t.mConnLock.Unlock()

	return t.mConnTransport.Endpoint()
}

func (t *InterceptedTransport) Accept(ctx context.Context) (Connection, error) {

	t.mConnLock.Lock()
	mConn := t.mConnTransport
	t.mConnLock.Unlock()
	c, err := mConn.Accept(ctx)
	if err != nil {
		return c, err
	}
	return NewInterceptedConnection(c, t.nodeID, t.netrixClient, t.router, t.logger), nil
}

func (t *InterceptedTransport) Dial(ctx context.Context, endpoint *Endpoint) (Connection, error) {
	t.mConnLock.Lock()
	mConn := t.mConnTransport
	t.mConnLock.Unlock()

	c, err := mConn.Dial(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	return NewInterceptedConnection(c, t.nodeID, t.netrixClient, t.router, t.logger), nil
}

func (t *InterceptedTransport) Listen(endpoint *Endpoint) error {
	t.mConnLock.Lock()
	mConn := t.mConnTransport
	t.mConnLock.Unlock()

	return mConn.Listen(endpoint)
}

func (t *InterceptedTransport) Close() error {
	t.mConnLock.Lock()
	mConn := t.mConnTransport
	t.mConnLock.Unlock()

	return mConn.Close()
}

func (t *InterceptedTransport) OnStart(context.Context) error {
	t.netrixClient.Start()
	go t.poll()
	return nil
}

func (t *InterceptedTransport) Ready() {
	t.netrixClient.Ready()
}

func (t *InterceptedTransport) NotReady() {
	t.netrixClient.NotReady()
}

func (t *InterceptedTransport) poll() {
	for {
		select {
		case <-t.doneCh:
		default:
			msg, ok := t.netrixClient.ReceiveMessage()
			if ok {
				go t.router.dispatchMessage(msg)
			}
		}
	}
}

func (t *InterceptedTransport) OnStop() {
	t.netrixClient.Stop()
	t.once.Do(func() {
		close(t.doneCh)
	})
}

func (t *InterceptedTransport) Reset() {
	t.mConnLock.Lock()
	t.mConnTransport = NewMConnTransport(
		t.logger,
		t.mConnConfig,
		t.channelDescs,
		t.mConnOptions,
	)
	t.mConnLock.Unlock()
}

func (t *InterceptedTransport) AddChannelDescriptors(d []*ChannelDescriptor) {
	t.mConnLock.Lock()
	mConn := t.mConnTransport
	t.mConnLock.Unlock()

	t.channelDescs = append(t.channelDescs, d...)
	mConn.AddChannelDescriptors(d)
}

type InterceptedConnection struct {
	conn         Connection
	senderID     string
	netrixClient *netrix.ReplicaClient
	recvCh       chan *ntypes.Message
	router       *router
	nodeID       types.NodeID
	doneCh       chan struct{}
	once         *sync.Once
	logger       log.Logger
}

var _ Connection = &InterceptedConnection{}

func NewInterceptedConnection(
	conn Connection,
	senderID string,
	client *netrix.ReplicaClient,
	router *router,
	logger log.Logger,
) *InterceptedConnection {
	return &InterceptedConnection{
		conn:         conn,
		senderID:     senderID,
		netrixClient: client,
		router:       router,
		recvCh:       make(chan *ntypes.Message, 10),
		doneCh:       make(chan struct{}),
		once:         new(sync.Once),
		logger:       logger,
	}
}

func (c *InterceptedConnection) Handshake(
	ctx context.Context,
	timeout time.Duration,
	nodeInfo types.NodeInfo,
	privKey crypto.PrivKey,
) (types.NodeInfo, crypto.PubKey, error) {
	nodeInfo, pubKey, err := c.conn.Handshake(ctx, timeout, nodeInfo, privKey)
	if err != nil {
		return nodeInfo, pubKey, err
	}
	c.nodeID = nodeInfo.NodeID
	c.router.addConnection(c.nodeID, c)
	return nodeInfo, pubKey, err
}

func (c *InterceptedConnection) ReceiveMessage(ctx context.Context) (ChannelID, []byte, error) {
	select {
	case m := <-c.recvCh:
		chanID, err := strconv.ParseUint(m.Type, 10, 16)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to parse message: %s", err)
		}
		return conn.ChannelID(chanID), m.Data, nil
	case <-c.doneCh:
		return 0, nil, io.EOF
	case <-ctx.Done():
		return 0, nil, io.EOF
	}
}

func (c *InterceptedConnection) SendMessage(ctx context.Context, cID ChannelID, data []byte) error {
	select {
	case <-ctx.Done():
		return io.EOF
	case <-c.doneCh:
		return io.EOF
	default:
		return c.netrixClient.SendMessageWithID(
			fmt.Sprintf("%s_%s_%s", c.senderID, c.nodeID, string(data)),
			strconv.FormatUint(uint64(cID), 10),
			ntypes.ReplicaID(c.nodeID),
			data,
		)
	}
}

func (c *InterceptedConnection) LocalEndpoint() Endpoint {
	return c.conn.LocalEndpoint()
}

func (c *InterceptedConnection) RemoteEndpoint() Endpoint {
	return c.conn.RemoteEndpoint()
}

func (c *InterceptedConnection) Close() error {
	c.once.Do(func() {
		close(c.doneCh)
	})
	return nil
}

func (c *InterceptedConnection) IsClosed() bool {
	select {
	case <-c.doneCh:
		return true
	default:
		return false
	}
}

func (c *InterceptedConnection) String() string {
	return c.conn.String()
}
