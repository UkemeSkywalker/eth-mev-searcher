package mempool

import (
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
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
}
