package mempool

import (
	beeUtils "github.com/astaxie/beego/utils"
	"github.com/btcboost/copernicus/algorithm"
	"github.com/btcboost/copernicus/model"
	"github.com/btcboost/copernicus/utils"
	"gopkg.in/fatih/set.v0"
)

/**
 * Fake height value used in Coins to signify they are only in the memory
 * pool(since 0.8)
 */
const (
	MEMPOOL_HEIGHT       = 0x7FFFFFFF
	ROLLING_FEE_HALFLIFE = 60 * 60 * 12
)

type Mempool struct {
	CheckFrequency              uint32
	TransactionsUpdated         int
	MinerPolicyEstimator        *BlockPolicyEstimator
	totalTxSize                 uint64
	CachedInnerUsage            uint64
	LastRollingFeeUpdate        int64
	BlockSinceLatRollingFeeBump bool
	RollingMinimumFeeRate       float64
	MapTx                       *beeUtils.BeeMap
	MapLinks                    *beeUtils.BeeMap //<TxMempoolEntry,Txlinks>
}

// UpdateForDescendants : Update the given tx for any in-mempool descendants.
// Assumes that setMemPoolChildren is correct for the given tx and all
// descendants.
func (mempool *Mempool) UpdateForDescendants(updateIt *TxMempoolEntry, cachedDescendants *algorithm.CacheMap, setExclude set.Set) {

	stageEntries := set.New()
	setAllDescendants := set.New()

	for !stageEntries.IsEmpty() {
		cit := stageEntries.List()[0]
		setAllDescendants.Add(cit)
		stageEntries.Remove(cit)
		txMempoolEntry := cit.(TxMempoolEntry)
		setChildren := mempool.GetMempoolChildren(&txMempoolEntry)

		for _, childEntry := range setChildren.Array {
			childTx := childEntry.(TxMempoolEntry)
			cacheIt := cachedDescendants.Get(childTx)
			cacheItVector := cacheIt.(algorithm.Vector)
			if cacheIt != cachedDescendants.Last() {
				// We've already calculated this one, just add the entries for
				// this set but don't traverse again.
				for _, cacheEntry := range cacheItVector.Array {
					setAllDescendants.Add(cacheEntry)

				}
			} else if !setAllDescendants.Has(childEntry) {
				// Schedule for later processing
				stageEntries.Add(childEntry)
			}

		}

	}
	// setAllDescendants now contains all in-mempool descendants of updateIt.
	// Update and add to cached descendant map
	modifySize := 0
	modifyFee := 0
	modifyCount := 0

	for _, cit := range setAllDescendants.List() {
		txCit := cit.(TxMempoolEntry)
		if !setExclude.Has(txCit.TxRef.Hash) {
			modifySize = modifySize + txCit.TxSize
			modifyFee = modifyFee + txCit.ModSize
			modifyCount++
			cachedSet := cachedDescendants.Get(updateIt).(set.Set)
			cachedSet.Add(txCit)
			// todo Update ancestor state for each descendant
		}
	}
	//todo Update descendant
}

func (mempool *Mempool) GetMempoolChildren(entry *TxMempoolEntry) *algorithm.Vector {
	result := mempool.MapLinks.Get(entry)
	if result == nil {
		panic("No have children In mempool for this TxmempoolEntry")
	}
	return result.(TxLinks).Children
}

func (mempool *Mempool) GetMemPoolParents(entry *TxMempoolEntry) *algorithm.Vector {
	result := mempool.MapLinks.Get(entry)
	if result == nil {
		panic("No have parant In mempool for this TxmempoolEntry")
	}
	return result.(TxLinks).Parents
}

func (mempool *Mempool) GetMinFee(sizeLimit uint) utils.FeeRate {
	return utils.FeeRate{SataoshisPerK: 0}
}

func AllowFee(priority float64) bool {
	// Large (in bytes) low-priority (new, small-coin) transactions need a fee.
	return priority > AllowFreeThreshold()
}

func GetTxFromMemPool(hash *utils.Hash) *model.Tx {
	return new(model.Tx)
}

func AllowFreeThreshold() float64 {
	return (float64(utils.COIN) * 144) / 250

}
