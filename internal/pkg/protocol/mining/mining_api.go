package mining

import (
	"context"

	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-filecoin/internal/pkg/block"
	"github.com/filecoin-project/go-filecoin/internal/pkg/clock"
	"github.com/filecoin-project/go-filecoin/internal/pkg/mining"
	"github.com/pkg/errors"
)

type miningChainReader interface {
	GetHead() block.TipSetKey
	GetTipSet(tsKey block.TipSetKey) (block.TipSet, error)
}

// API provides an interface to the block mining protocol.
type API struct {
	minerAddress    func() (address.Address, error)
	addNewBlockFunc func(context.Context, *block.Block) (err error)
	chainReader     miningChainReader
	isMiningFunc    func() bool
	setupMiningFunc func(context.Context) error
	startMiningFunc func(context.Context) error
	stopMiningFunc  func(context.Context)
	getWorkerFunc   func(ctx context.Context) (*mining.DefaultWorker, error)
	chainClock      clock.ChainEpochClock
}

// New creates a new API instance with the provided deps
func New(
	minerAddr func() (address.Address, error),
	addNewBlockFunc func(context.Context, *block.Block) (err error),
	chainReader miningChainReader,
	isMiningFunc func() bool,
	setupMiningFunc func(ctx context.Context) error,
	startMiningFunc func(context.Context) error,
	stopMiningfunc func(context.Context),
	getWorkerFunc func(ctx context.Context) (*mining.DefaultWorker, error),
	chainClock clock.ChainEpochClock,
) API {
	return API{
		minerAddress:    minerAddr,
		addNewBlockFunc: addNewBlockFunc,
		chainReader:     chainReader,
		isMiningFunc:    isMiningFunc,
		setupMiningFunc: setupMiningFunc,
		startMiningFunc: startMiningFunc,
		stopMiningFunc:  stopMiningfunc,
		getWorkerFunc:   getWorkerFunc,
		chainClock:      chainClock,
	}
}

// MinerAddress returns the mining address the API is using, an error is
// returned if the mining address is not set.
func (a *API) MinerAddress() (address.Address, error) {
	return a.minerAddress()
}

// MiningIsActive calls the node's IsMining function
func (a *API) MiningIsActive() bool {
	return a.isMiningFunc()
}

// MiningOnce mines a single block in the given context, and returns the new block.
func (a *API) MiningOnce(ctx context.Context) (*block.Block, error) {
	if a.isMiningFunc() {
		return nil, errors.New("Node is already mining")
	}

	ts, err := a.chainReader.GetTipSet(a.chainReader.GetHead())
	if err != nil {
		return nil, err
	}

	miningWorker, err := a.getWorkerFunc(ctx)
	if err != nil {
		return nil, err
	}

	res, err := mining.MineOnce(ctx, *miningWorker, ts, a.chainClock)
	if err != nil {
		return nil, err
	}
	if res.Err != nil {
		return nil, res.Err
	}

	if err := a.addNewBlockFunc(ctx, res.NewBlock); err != nil {
		return nil, err
	}

	return res.NewBlock, nil
}

// MiningSetup sets up a storage miner without running repeated tasks like mining
func (a *API) MiningSetup(ctx context.Context) error {
	return a.setupMiningFunc(ctx)
}

// MiningStart calls the node's StartMining function
func (a *API) MiningStart(ctx context.Context) error {
	return a.startMiningFunc(ctx)
}

// MiningStop calls the node's StopMining function
func (a *API) MiningStop(ctx context.Context) {
	a.stopMiningFunc(ctx)
}
