package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	netrix "github.com/netrixframework/go-clientlibrary"
	ntypes "github.com/netrixframework/netrix/types"
	abciclient "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/internal/blocksync"
	"github.com/tendermint/tendermint/internal/consensus"
	"github.com/tendermint/tendermint/internal/eventbus"
	"github.com/tendermint/tendermint/internal/eventlog"
	"github.com/tendermint/tendermint/internal/p2p"
	"github.com/tendermint/tendermint/internal/p2p/conn"
	"github.com/tendermint/tendermint/internal/p2p/pex"
	"github.com/tendermint/tendermint/internal/proxy"
	rpccore "github.com/tendermint/tendermint/internal/rpc/core"
	sm "github.com/tendermint/tendermint/internal/state"
	"github.com/tendermint/tendermint/internal/state/indexer"
	"github.com/tendermint/tendermint/internal/state/indexer/sink"
	"github.com/tendermint/tendermint/internal/statesync"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/service"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

func NetrixNode(
	ctx context.Context,
	conf *config.Config,
	logger log.Logger,
) (service.Service, error) {
	return newNetrixNode(ctx, conf, logger)
}

type netrixNode struct {
	node        *nodeImpl
	nodeCtx     context.CancelFunc
	filePrivVal *privval.FilePV
	config      *config.Config
	logger      log.Logger

	transport *p2p.InterceptedTransport
	lock      *sync.Mutex
}

type netrixNodeServiceWrapper struct {
	node *netrixNode
	service.BaseService
}

func (n *netrixNodeServiceWrapper) OnStart(context.Context) error {
	n.node.Start()
	return nil
}

func (n *netrixNodeServiceWrapper) OnStop() {
	n.node.Stop()
}

func nodeInitalSetup(cfg *config.Config) (types.NodeKey, *privval.FilePV, *types.GenesisDoc, error) {
	nodeKey, err := types.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return types.NodeKey{}, nil, nil, fmt.Errorf("failed to load or gen node key %s: %w", cfg.NodeKeyFile(), err)
	}

	filePrivval, err := makeDefaultPrivval(cfg)
	if err != nil {
		return types.NodeKey{}, nil, nil, err
	}

	genesisDocProvider := defaultGenesisDocProviderFunc(cfg)

	genDoc, err := genesisDocProvider()
	if err != nil {
		return types.NodeKey{}, nil, nil, fmt.Errorf("failed to read genesis file: %s", err)
	}

	if err = genDoc.ValidateAndComplete(); err != nil {
		return types.NodeKey{}, nil, nil, fmt.Errorf("error in genesis doc: %w", err)
	}
	return nodeKey, filePrivval, genDoc, nil
}

func newNetrixNode(ctx context.Context, cfg *config.Config, logger log.Logger) (service.Service, error) {

	nodeKey, filePrivval, genDoc, err := nodeInitalSetup(cfg)
	if err != nil {
		return nil, err
	}
	if filePrivval == nil {
		return nil, fmt.Errorf("no key")
	}
	if genDoc == nil {
		return nil, fmt.Errorf("no gen doc")
	}

	node := &netrixNode{
		filePrivVal: filePrivval,
		config:      cfg,
		logger:      logger,
		lock:        new(sync.Mutex),
	}

	transport, err := createNetrixTransport(logger, cfg.P2P, node, nodeKey.ID, genDoc.ChainID, &filePrivval.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %s", err)
	}
	node.transport = transport

	wrapper := &netrixNodeServiceWrapper{
		node: node,
	}
	wrapper.BaseService = *service.NewBaseService(logger.With("service", "netrix-node"), "netrix-node", wrapper)
	return wrapper, nil
}

func createRouterWithTransport(
	logger log.Logger,
	p2pMetrics *p2p.Metrics,
	nodeInfoProducer func() *types.NodeInfo,
	nodeKey types.NodeKey,
	peerManager *p2p.PeerManager,
	cfg *config.Config,
	appClient abciclient.Client,
	transport p2p.Transport,
) (*p2p.Router, error) {
	p2pLogger := logger.With("module", "p2p")
	ep, err := p2p.NewEndpoint(nodeKey.ID.AddressString(cfg.P2P.ListenAddress))
	if err != nil {
		return nil, err
	}

	return p2p.NewRouter(
		p2pLogger,
		p2pMetrics,
		nodeKey.PrivKey,
		peerManager,
		nodeInfoProducer,
		transport,
		ep,
		getRouterConfig(cfg, appClient),
	)
}

func createNetrixTransport(
	logger log.Logger,
	p2pConfig *config.P2PConfig,
	dh netrix.DirectiveHandler,
	nodeID types.NodeID,
	chainID string,
	pvKey *privval.FilePVKey,
) (*p2p.InterceptedTransport, error) {

	keyS, err := json.Marshal(pvKey)
	if err != nil {
		keyS = []byte{}
	}

	transportConf := conn.DefaultMConnConfig()
	transportConf.FlushThrottle = p2pConfig.FlushThrottleTimeout
	transportConf.SendRate = p2pConfig.SendRate
	transportConf.RecvRate = p2pConfig.RecvRate
	transportConf.MaxPacketMsgPayloadSize = p2pConfig.MaxPacketMsgPayloadSize

	transport, err := p2p.NewInterceptedTransport(
		logger,
		transportConf,
		[]*p2p.ChannelDescriptor{},
		p2p.InterceptedTransportOptions{
			MConnOptions: p2p.MConnTransportOptions{
				MaxAcceptedConnections: uint32(p2pConfig.MaxConnections),
			},
			NetrixConfig: &netrix.Config{
				ReplicaID:        ntypes.ReplicaID(nodeID),
				ClientServerAddr: p2pConfig.NetrixClientAddr,
				ClientAdvAddr:    p2pConfig.NetrixClientAdvAddr,
				NetrixAddr:       p2pConfig.NetrixAddr,
				Info: map[string]interface{}{
					"chain_id": chainID,
					"privKey":  string(keyS),
				},
			},
		},
		dh,
	)
	return transport, err
}

func (n *netrixNode) createNode() (*nodeImpl, error) {

	cfg := n.config

	nodeKey, _, genDoc, err := nodeInitalSetup(cfg)
	if err != nil {
		return nil, err
	}
	n.logger.With("node_key", nodeKey).Info("Fetched node key!")
	filePrivval := n.filePrivVal

	dbProvider := config.DefaultDBProvider

	var cancel context.CancelFunc
	ctx, cancel := context.WithCancel(context.Background())

	closers := []closer{convertCancelCloser(cancel)}

	blockStore, stateDB, dbCloser, err := initDBs(cfg, dbProvider)
	if err != nil {
		return nil, combineCloseError(err, dbCloser)
	}
	closers = append(closers, dbCloser)

	stateStore := sm.NewStore(stateDB)

	state, err := loadStateFromDBOrGenesisDocProvider(stateStore, genDoc)
	if err != nil {
		return nil, combineCloseError(err, makeCloser(closers))
	}

	client, _, err := proxy.ClientFactory(n.logger, cfg.ProxyApp, cfg.ABCI, cfg.DBDir())
	if err != nil {
		return nil, err
	}

	nodeMetrics := defaultMetricsProvider(cfg.Instrumentation)(genDoc.ChainID)

	proxyApp := proxy.New(client, n.logger.With("module", "proxy"), nodeMetrics.proxy)
	eventBus := eventbus.NewDefault(n.logger.With("module", "events"))

	var eventLog *eventlog.Log
	if w := cfg.RPC.EventLogWindowSize; w > 0 {
		var err error
		eventLog, err = eventlog.New(eventlog.LogSettings{
			WindowSize: w,
			MaxItems:   cfg.RPC.EventLogMaxItems,
			Metrics:    nodeMetrics.eventlog,
		})
		if err != nil {
			return nil, combineCloseError(fmt.Errorf("initializing event log: %w", err), makeCloser(closers))
		}
	}
	eventSinks, err := sink.EventSinksFromConfig(cfg, dbProvider, genDoc.ChainID)
	if err != nil {
		return nil, combineCloseError(err, makeCloser(closers))
	}
	indexerService := indexer.NewService(indexer.ServiceArgs{
		Sinks:    eventSinks,
		EventBus: eventBus,
		Logger:   n.logger.With("module", "txindex"),
		Metrics:  nodeMetrics.indexer,
	})

	privValidator, err := createPrivval(ctx, n.logger, cfg, genDoc, filePrivval)
	if err != nil {
		return nil, combineCloseError(err, makeCloser(closers))
	}

	var pubKey crypto.PubKey
	if cfg.Mode == config.ModeValidator {
		pubKey, err = privValidator.GetPubKey(ctx)
		if err != nil {
			return nil, combineCloseError(fmt.Errorf("can't get pubkey: %w", err),
				makeCloser(closers))

		}
		if pubKey == nil {
			return nil, combineCloseError(
				errors.New("could not retrieve public key from private validator"),
				makeCloser(closers))
		}
	}

	peerManager, peerCloser, err := createPeerManager(cfg, dbProvider, nodeKey.ID, nodeMetrics.p2p)
	closers = append(closers, peerCloser)
	if err != nil {
		return nil, combineCloseError(
			fmt.Errorf("failed to create peer manager: %w", err),
			makeCloser(closers))
	}

	// TODO construct node here:
	node := &nodeImpl{
		config:        cfg,
		logger:        n.logger,
		genesisDoc:    genDoc,
		privValidator: privValidator,

		peerManager: peerManager,
		nodeKey:     nodeKey,

		eventSinks:     eventSinks,
		indexerService: indexerService,
		services:       []service.Service{eventBus},

		initialState: state,
		stateStore:   stateStore,
		blockStore:   blockStore,

		shutdownOps: makeCloser(closers),

		rpcEnv: &rpccore.Environment{
			ProxyApp: proxyApp,

			StateStore: stateStore,
			BlockStore: blockStore,

			PeerManager: peerManager,

			GenDoc:     genDoc,
			EventSinks: eventSinks,
			EventBus:   eventBus,
			EventLog:   eventLog,
			Logger:     n.logger.With("module", "rpc"),
			Config:     *cfg.RPC,
		},
	}

	node.router, err = createRouterWithTransport(n.logger, nodeMetrics.p2p, node.NodeInfo, nodeKey, peerManager, cfg, proxyApp, n.transport)
	if err != nil {
		return nil, combineCloseError(
			fmt.Errorf("failed to create router: %w", err),
			makeCloser(closers))
	}

	evReactor, evPool, edbCloser, err := createEvidenceReactor(n.logger, cfg, dbProvider,
		stateStore, blockStore, peerManager.Subscribe, node.router.OpenChannel, nodeMetrics.evidence, eventBus)
	closers = append(closers, edbCloser)
	if err != nil {
		return nil, combineCloseError(err, makeCloser(closers))
	}
	node.services = append(node.services, evReactor)
	node.rpcEnv.EvidencePool = evPool
	node.evPool = evPool

	mpReactor, mp := createMempoolReactor(n.logger, cfg, proxyApp, stateStore, nodeMetrics.mempool,
		peerManager.Subscribe, node.router.OpenChannel)
	node.rpcEnv.Mempool = mp
	node.services = append(node.services, mpReactor)

	// make block executor for consensus and blockchain reactors to execute blocks
	blockExec := sm.NewBlockExecutor(
		stateStore,
		n.logger.With("module", "state"),
		proxyApp,
		mp,
		evPool,
		blockStore,
		eventBus,
		nodeMetrics.state,
	)

	// Determine whether we should attempt state sync.
	stateSync := cfg.StateSync.Enable && !onlyValidatorIsUs(state, pubKey)
	if stateSync && state.LastBlockHeight > 0 {
		n.logger.Info("Found local state with non-zero height, skipping state sync")
		stateSync = false
	}

	// Determine whether we should do block sync. This must happen after the handshake, since the
	// app may modify the validator set, specifying ourself as the only validator.
	blockSync := !onlyValidatorIsUs(state, pubKey)
	waitSync := stateSync || blockSync

	var ticker consensus.TimeoutTicker
	if cfg.Consensus.NetrixTimer {
		ticker, err = consensus.NewNetrixTicker(n.logger.With("module", "netrix-ticker"))
		if err != nil {
			return nil, combineCloseError(err, makeCloser(closers))
		}
	} else {
		ticker = consensus.NewTimeoutTicker(n.logger)
	}

	csState, err := consensus.NewState(n.logger.With("module", "consensus"),
		cfg.Consensus,
		stateStore,
		blockExec,
		blockStore,
		mp,
		ticker,
		evPool,
		eventBus,
		consensus.StateMetrics(nodeMetrics.consensus),
		consensus.SkipStateStoreBootstrap,
	)
	if err != nil {
		return nil, combineCloseError(err, makeCloser(closers))
	}
	node.rpcEnv.ConsensusState = csState

	csReactor := consensus.NewReactor(
		n.logger,
		csState,
		node.router.OpenChannel,
		peerManager.Subscribe,
		eventBus,
		waitSync,
		nodeMetrics.consensus,
	)
	node.services = append(node.services, csReactor)
	node.rpcEnv.ConsensusReactor = csReactor

	// Create the blockchain reactor. Note, we do not start block sync if we're
	// doing a state sync first.
	bcReactor := blocksync.NewReactor(
		n.logger.With("module", "blockchain"),
		stateStore,
		blockExec,
		blockStore,
		csReactor,
		node.router.OpenChannel,
		peerManager.Subscribe,
		blockSync && !stateSync,
		nodeMetrics.consensus,
		eventBus,
	)
	node.services = append(node.services, bcReactor)
	node.rpcEnv.BlockSyncReactor = bcReactor

	// Make ConsensusReactor. Don't enable fully if doing a state sync and/or block sync first.
	// FIXME We need to update metrics here, since other reactors don't have access to them.
	if stateSync {
		nodeMetrics.consensus.StateSyncing.Set(1)
	} else if blockSync {
		nodeMetrics.consensus.BlockSyncing.Set(1)
	}

	if cfg.P2P.PexReactor {
		node.services = append(node.services, pex.NewReactor(n.logger, peerManager, node.router.OpenChannel, peerManager.Subscribe))
	}

	// Set up state sync reactor, and schedule a sync if requested.
	// FIXME The way we do phased startups (e.g. replay -> block sync -> consensus) is very messy,
	// we should clean this whole thing up. See:
	// https://github.com/tendermint/tendermint/issues/4644
	node.services = append(node.services, statesync.NewReactor(
		genDoc.ChainID,
		genDoc.InitialHeight,
		*cfg.StateSync,
		n.logger.With("module", "statesync"),
		proxyApp,
		node.router.OpenChannel,
		peerManager.Subscribe,
		stateStore,
		blockStore,
		cfg.StateSync.TempDir,
		nodeMetrics.statesync,
		eventBus,
		// the post-sync operation
		func(ctx context.Context, state sm.State) error {
			csReactor.SetStateSyncingMetrics(0)

			// TODO: Some form of orchestrator is needed here between the state
			// advancing reactors to be able to control which one of the three
			// is running
			// FIXME Very ugly to have these metrics bleed through here.
			csReactor.SetBlockSyncingMetrics(1)
			if err := bcReactor.SwitchToBlockSync(ctx, state); err != nil {
				n.logger.Error("failed to switch to block sync", "err", err)
				return err
			}

			return nil
		},
		stateSync,
	))

	if cfg.Mode == config.ModeValidator {
		if privValidator != nil {
			csState.SetPrivValidator(ctx, privValidator)
		}
		node.rpcEnv.PubKey = pubKey
	}

	node.BaseService = *service.NewBaseService(n.logger, "Node", node)

	return node, nil
}

func (n *netrixNode) Start() error {
	if !n.transport.IsRunning() {
		n.transport.Start(context.Background())
	}

	n.lock.Lock()
	node := n.node
	n.lock.Unlock()
	if node == nil {
		n.logger.Info("Creating new node")
		newNode, err := n.createNode()
		if err != nil {
			return err
		}
		n.lock.Lock()
		n.node = newNode
		n.lock.Unlock()
	}
	n.logger.Info("Starting node")
	ctx, cancel := context.WithCancel(context.Background())
	n.lock.Lock()
	node = n.node
	n.nodeCtx = cancel
	n.lock.Unlock()
	err := node.Start(ctx)
	if err != nil {
		return err
	}
	n.logger.Info("Started node")
	n.transport.Ready()
	return nil
}

func (n *netrixNode) Stop() error {
	n.lock.Lock()
	if n.node != nil {
		n.nodeCtx()
		n.node.Wait()
		n.logger.Info("Stopped node successfully")
	}
	n.lock.Unlock()
	return nil
}

func (n *netrixNode) Restart() error {
	n.transport.NotReady()
	n.logger.Info("Stopping node")
	n.Stop()
	n.logger.Info("Clearing storage")
	n.clearStorage()
	n.filePrivVal.Reset()
	n.lock.Lock()
	n.node = nil
	n.nodeCtx = nil
	n.lock.Unlock()
	n.transport.Reset()
	return n.Start()
}

func (n *netrixNode) clearStorage() error {
	os.RemoveAll(n.config.DBDir())
	os.MkdirAll(n.config.DBDir(), 0755)
	return nil
}
