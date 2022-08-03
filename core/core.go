package core

import (
	"math/big"

	"github.com/spruce-solutions/go-quai/common"
	"github.com/spruce-solutions/go-quai/consensus"
	"github.com/spruce-solutions/go-quai/core/state"
	"github.com/spruce-solutions/go-quai/core/types"
	"github.com/spruce-solutions/go-quai/core/vm"
	"github.com/spruce-solutions/go-quai/ethdb"
	"github.com/spruce-solutions/go-quai/event"
	"github.com/spruce-solutions/go-quai/params"
)

type Core struct {
	sl     *Slice
	engine consensus.Engine
}

func NewCore(db ethdb.Database, chainConfig *params.ChainConfig, domClientUrl string, subClientUrls []string, engine consensus.Engine, vmConfig vm.Config) (*Core, error) {

	slice, err := NewSlice(db, chainConfig, domClientUrl, subClientUrls, engine, vmConfig)
	if err != nil {
		return nil, err
	}

	return &Core{
		sl:     slice,
		engine: engine,
	}, nil
}

func (c *Core) Config() *params.ChainConfig {
	return c.sl.hc.bc.chainConfig
}

// Engine retreives the blake3 consensus engine.
func (c *Core) Engine() consensus.Engine {
	return c.engine
}

// Slice retrieves the slice struct.
func (c *Core) Slice() *Slice {
	return c.sl
}

// GetBlock retrieves a block from the database by hash and number,
// caching it if found.
func (c *Core) GetBlock(hash common.Hash, number uint64) *types.Block {
	return c.sl.hc.GetBlock(hash, number)
}

// GetBlockByHash retrieves a block from the database by hash, caching it if found.
func (c *Core) GetBlockByHash(hash common.Hash) *types.Block {
	return c.sl.hc.GetBlockByHash(hash)
}

// GetHeaderByNumber retrieves a block header from the database by number,
// caching it (associated with its hash) if found.
func (c *Core) GetHeaderByNumber(number uint64) *types.Header {
	return c.sl.hc.GetHeaderByNumber(number)
}

// GetBlockByNumber retrieves a block from the database by number, caching it
// (associated with its hash) if found.
func (c *Core) GetBlockByNumber(number uint64) *types.Block {
	return c.sl.hc.GetBlockByNumber(number)
}

// GetBlocksFromHash returns the block corresponding to hash and up to n-1 ancestors.
// [deprecated by eth/62]
func (c *Core) GetBlocksFromHash(hash common.Hash, n int) []*types.Block {
	return c.sl.hc.GetBlocksFromHash(hash, n)
}

// GetUnclesInChain retrieves all the uncles from a given block backwards until
// a specific distance is reached.
func (c *Core) GetUnclesInChain(block *types.Block, length int) []*types.Header {
	return c.sl.hc.GetUnclesInChain(block, length)
}

// GetGasUsedInChain retrieves all the gas used from a given block backwards until
// a specific distance is reached.
func (c *Core) GetGasUsedInChain(block *types.Block, length int) int64 {
	return c.sl.hc.GetGasUsedInChain(block, length)
}

// GetGasUsedInChain retrieves all the gas used from a given block backwards until
// a specific distance is reached.
func (c *Core) CalculateBaseFee(header *types.Header) *big.Int {
	return c.sl.hc.CalculateBaseFee(header)
}

// CurrentBlock returns the block for the current header.
func (c *Core) CurrentBlock() *types.Block {
	return c.GetBlockByHash(c.CurrentHeader().Hash())
}

// CurrentHeader retrieves the current head header of the canonical chain. The
// header is retrieved from the HeaderChain's internal cache.
func (c *Core) CurrentHeader() *types.Header {
	return c.sl.hc.CurrentHeader()
}

// GetTd retrieves a block's total difficulty in the canonical chain from the
// database by hash and number, caching it if found.
func (c *Core) GetTd(hash common.Hash, number uint64) []*big.Int {
	return c.sl.hc.GetTd(hash, number)
}

// GetTdByHash retrieves a block's total difficulty in the canonical chain from the
// database by hash, caching it if found.
func (c *Core) GetTdByHash(hash common.Hash) []*big.Int {
	return c.sl.hc.GetTdByHash(hash)
}

// GetHeader retrieves a block header from the database by hash and number,
// caching it if found.
func (c *Core) GetHeader(hash common.Hash, number uint64) *types.Header {
	return c.sl.hc.GetHeader(hash, number)
}

// GetHeaderByHash retrieves a block header from the database by hash, caching it if
// found.
func (c *Core) GetHeaderByHash(hash common.Hash) *types.Header {
	return c.sl.hc.GetHeaderByHash(hash)
}

// HasBlock checks if a block is fully present in the database or not.
func (c *Core) HasBlock(hash common.Hash, number uint64) bool {
	return c.sl.hc.bc.HasBlock(hash, number)
}

// HasHeader checks if a block header is present in the database or not, caching
// it if present.
func (c *Core) HasHeader(hash common.Hash, number uint64) bool {
	return c.sl.hc.HasHeader(hash, number)
}

// GetCanonicalHash returns the canonical hash for a given block number
func (c *Core) GetCanonicalHash(number uint64) common.Hash {
	return c.sl.hc.GetCanonicalHash(number)
}

// GetBlockHashesFromHash retrieves a number of block hashes starting at a given
// hash, fetching towards the genesis block.
func (c *Core) GetBlockHashesFromHash(hash common.Hash, max uint64) []common.Hash {
	return c.sl.hc.GetBlockHashesFromHash(hash, max)
}

// GetAncestor retrieves the Nth ancestor of a given block. It assumes that either the given block or
// a close ancestor of it is canonical. maxNonCanonical points to a downwards counter limiting the
// number of blocks to be individually checked before we reach the canonical chain.
//
// Note: ancestor == 0 returns the same block, 1 returns its parent and so on.
func (c *Core) GetAncestor(hash common.Hash, number, ancestor uint64, maxNonCanonical *uint64) (common.Hash, uint64) {
	return c.sl.hc.GetAncestor(hash, number, ancestor, maxNonCanonical)
}

// GetAncestorWithLocation retrieves the first occurrence of a block with a given location from a given block.
//
// Note: location == hash location returns the same block.
func (c *Core) GetAncestorByLocation(hash common.Hash, location []byte) (*types.Header, error) {
	return c.sl.hc.GetAncestorByLocation(hash, location)
}

// ContractCode retrieves a blob of data associated with a contract hash
// either from ephemeral in-memory cache, or from persistent storage.
func (c *Core) ContractCode(hash common.Hash) ([]byte, error) {
	return c.sl.hc.bc.processor.ContractCode(hash)
}

// StateAt returns a new mutable state based on a particular point in time.
func (c *Core) StateAt(root common.Hash) (*state.StateDB, error) {
	return state.New(root, c.sl.hc.bc.processor.stateCache, nil)
}

// ContractCodeWithPrefix retrieves a blob of data associated with a contract
// hash either from ephemeral in-memory cache, or from persistent storage.
//
// If the code doesn't exist in the in-memory cache, check the storage with
// new code scheme.
func (c *Core) ContractCodeWithPrefix(hash common.Hash) ([]byte, error) {
	return c.sl.hc.bc.processor.ContractCodeWithPrefix(hash)
}

// Genesis retrieves the chain's genesis block.
func (c *Core) Genesis() *types.Block {
	return c.GetBlockByHash(c.sl.hc.genesisHeader.Hash())
}
func (c *Core) ResetWithGenesisBlock(genesis *types.Header) error {
	return c.sl.hc.ResetWithGenesisBlock(genesis)
}

func (c *Core) Stop() {
	c.sl.hc.bc.Stop()
}

// SubscribeChainHeadEvent registers a subscription of ChainHeadEvent.
func (c *Core) SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription {
	return c.sl.hc.scope.Track(c.sl.hc.chainHeadFeed.Subscribe(ch))
}

// GetDifficultyOrder determines the difficulty order of the given header.
func (c *Core) GetDifficultyOrder(header *types.Header) (int, error) {
	return c.sl.GetDifficultyOrder(header)
}

// HLCR does hierarchical comparison of two difficulty tuples and returns true if second tuple is greater than the first
func (c *Core) HLCR(localDifficulties []*big.Int, externDifficulties []*big.Int) bool {
	return c.sl.HLCR(localDifficulties, externDifficulties)
}