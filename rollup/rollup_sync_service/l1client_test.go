package rollup_sync_service

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/scroll-tech/go-ethereum"
	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/core/types"
	"github.com/scroll-tech/go-ethereum/rlp"
)

func TestL1Client(t *testing.T) {
	ctx := context.Background()
	mockClient := &mockEthClient{}

	scrollChainABI, err := ScrollChainMetaData.GetAbi()
	if err != nil {
		t.Fatal("failed to get scroll chain abi", "err", err)
	}
	scrollChainAddress := common.HexToAddress("0x0123456789abcdef")
	l1Client, err := NewL1Client(ctx, mockClient, 11155111, scrollChainAddress, scrollChainABI)
	require.NoError(t, err, "Failed to initialize L1Client")

	blockNumber, err := l1Client.GetLatestFinalizedBlockNumber()
	assert.NoError(t, err, "Error getting latest confirmed block number")
	assert.Equal(t, uint64(36), blockNumber, "Unexpected block number")

	logs, err := l1Client.FetchRollupEventsInRange(0, blockNumber)
	assert.NoError(t, err, "Error fetching rollup events in range")
	assert.Empty(t, logs, "Expected no logs from FetchRollupEventsInRange")
}

type mockEthClient struct {
	txRLP []byte
}

func (m *mockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	return 11155111, nil
}

func (m *mockEthClient) ChainID(ctx context.Context) (*big.Int, error) {
	return big.NewInt(11155111), nil
}

func (m *mockEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return []types.Log{}, nil
}

func (m *mockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return &types.Header{
		Number: big.NewInt(100 - 64),
	}, nil
}

func (m *mockEthClient) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return nil, nil
}

func (m *mockEthClient) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	var tx types.Transaction
	if err := rlp.DecodeBytes(m.txRLP, &tx); err != nil {
		return nil, false, err
	}
	return &tx, false, nil
}

func (m *mockEthClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return nil, nil
}
