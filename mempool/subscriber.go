package mempool

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Subscriber struct {
	wsURL   string
	handler *Handler
}

func NewSubscriber(wsURL string) *Subscriber {
	return &Subscriber{
		wsURL:   wsURL,
		handler: NewHandler(),
	}
}

func (s *Subscriber) Start(ctx context.Context) error {
	for {
		if err := s.run(ctx); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("Connection error: %v, reconnecting in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return nil
	}
}

func (s *Subscriber) run(ctx context.Context) error {
	rpcClient, err := rpc.DialContext(ctx, s.wsURL)
	if err != nil {
		return err
	}
	defer rpcClient.Close()

	client := ethclient.NewClient(rpcClient)
	defer client.Close()

	log.Println("Connected to Geth, subscribing to pending transactions...")

	pendingTxHashCh := make(chan common.Hash, 1000)
	sub, err := rpcClient.EthSubscribe(ctx, pendingTxHashCh, "newPendingTransactions")
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	log.Println("Subscribed to mempool")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			return err
		case txHash := <-pendingTxHashCh:
			go s.fetchAndHandle(ctx, client, txHash)
		}
	}
}

func (s *Subscriber) fetchAndHandle(ctx context.Context, client *ethclient.Client, txHash common.Hash) {
	tx, _, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		log.Printf("Failed to fetch tx %s: %v", txHash.Hex(), err)
		return
	}
	s.handler.Handle(tx)
}
