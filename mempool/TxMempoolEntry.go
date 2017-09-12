package mempool

import "github.com/btcboost/copernicus/model"

type TxMempoolEntry struct {
	TxRef         *model.Tx
	Fee           int64
	TxSize        int
	UsageSize     int
	LocalTime     int64
	EntryPriority float64
	EntryHeight   int
	//!< Sum of all txin values that are already in blockchain
	InChainInputValue int64
	SpendsCoinbase    bool
	SigOpCount        int64
	FeeDelta          int64

	LockPoints *LockPoints
	// Information about descendants of this transaction that are in the
	// mempool; if we remove this transaction we must remove all of these
	// descendants as well.  if nCountWithDescendants is 0, treat this entry as
	// dirty, and nSizeWithDescendants and nModFeesWithDescendants will not be
	// correct.
	//!< number of descendant transactions
	CountWithDescendants    uint64
	SizeWithDescendants     uint64
	ModFeesWithDescendants  int64
	SigOpCoungWithAncestors int64
}