package mempool

import (
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/flashbots-lab/searcher/bundle"
)

type Handler struct{
	bundler *bundle.Bundler
}

func NewHandler(bundler *bundle.Bundler) *Handler {
	return &Handler{
		bundler: bundler,
	}
}

func (h *Handler) Handle(tx *types.Transaction) {
	to := "nil"
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	log.Printf("TX | hash=%s to=%s value=%s gas=%d input_len=%d",
		tx.Hash().Hex(),
		to,
		tx.Value().String(),
		tx.Gas(),
		len(tx.Data()),
	)

	h.bundler.Add(tx)
}
