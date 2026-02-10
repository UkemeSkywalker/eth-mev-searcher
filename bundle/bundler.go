package bundle

import (
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

const (
	MaxTxsPerBundle = 5
	BundleTimeout   = 10 * time.Second
)

type Bundle struct {
	Transactions []*types.Transaction
	TotalValue   *big.Int
	TotalGas     uint64
}

type Bundler struct {
	mu           sync.Mutex
	current      *Bundle
	timer        *time.Timer
	finalizeChan chan *Bundle
}

func NewBundler() *Bundler {
	b := &Bundler{
		current:      newBundle(),
		finalizeChan: make(chan *Bundle, 10),
	}
	b.resetTimer()
	return b
}

func newBundle() *Bundle {
	return &Bundle{
		Transactions: make([]*types.Transaction, 0, MaxTxsPerBundle),
		TotalValue:   big.NewInt(0),
		TotalGas:     0,
	}
}

func (b *Bundler) Add(tx *types.Transaction) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.current.Transactions = append(b.current.Transactions, tx)
	b.current.TotalValue.Add(b.current.TotalValue, tx.Value())
	b.current.TotalGas += tx.Gas()

	if len(b.current.Transactions) >= MaxTxsPerBundle {
		b.finalize()
	}
}

func (b *Bundler) finalize() {
	if len(b.current.Transactions) == 0 {
		return
	}

	b.timer.Stop()
	b.finalizeChan <- b.current
	b.current = newBundle()
	b.resetTimer()
}

func (b *Bundler) resetTimer() {
	b.timer = time.AfterFunc(BundleTimeout, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		b.finalize()
	})
}

func (b *Bundler) Start() <-chan *Bundle {
	return b.finalizeChan
}

func LogBundle(bundle *Bundle) {
	hashes := make([]string, len(bundle.Transactions))
	for i, tx := range bundle.Transactions {
		hashes[i] = tx.Hash().Hex()
	}

	log.Printf("BUNDLE | txs=%d total_value=%s total_gas=%d hashes=%v",
		len(bundle.Transactions),
		bundle.TotalValue.String(),
		bundle.TotalGas,
		hashes,
	)
}
