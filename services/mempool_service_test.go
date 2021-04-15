// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"context"
	"testing"

	"github.com/DeFiCh/rosetta-defichain/defichain"

	"github.com/DeFiCh/rosetta-defichain/configuration"
	mocks "github.com/DeFiCh/rosetta-defichain/mocks/services"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestMempoolEndpoints_Offline(t *testing.T) {
	cfg := &configuration.Configuration{
		Mode: configuration.Offline,
	}
	mockClient := &mocks.Client{}
	servicer := NewMempoolAPIService(cfg, mockClient)
	ctx := context.Background()
	mem, err := servicer.Mempool(ctx, nil)
	assert.Nil(t, mem)
	assert.Equal(t, ErrUnavailableOffline.Code, err.Code)
	assert.Equal(t, ErrUnavailableOffline.Message, err.Message)

	memTransaction, err := servicer.MempoolTransaction(ctx, nil)
	assert.Nil(t, memTransaction)
	assert.Equal(t, ErrUnavailableOffline.Code, err.Code)
	assert.Equal(t, ErrUnavailableOffline.Message, err.Message)
	mockClient.AssertExpectations(t)
}

func TestMempoolEndpoints_Online(t *testing.T) {
	cfg := &configuration.Configuration{
		Mode: configuration.Online,
	}

	mockClient := &mocks.Client{}
	servicer := NewMempoolAPIService(cfg, mockClient)
	ctx := context.Background()

	mockClient.On("RawMempool", ctx).Return([]string{
		"tx1",
		"tx2",
	}, nil)
	mem, err := servicer.Mempool(ctx, nil)
	assert.Nil(t, err)
	assert.Equal(t, &types.MempoolResponse{
		TransactionIdentifiers: []*types.TransactionIdentifier{
			{
				Hash: "tx1",
			},
			{
				Hash: "tx2",
			},
		},
	}, mem)

	mockClient.On("GetRawTransaction", ctx, "tx_hash", "").Return(&defichain.Transaction{
		Hash: "tx_hash",
		Inputs: []*defichain.Input{
			{
				TxHash: "input_tx_hash",
				Vout:   0,
			},
		},
		Outputs: []*defichain.Output{
			{
				Value: 1000000000,
				Index: 0,
			},
		},
	}, nil)
	memTransaction, err := servicer.MempoolTransaction(ctx, &types.MempoolTransactionRequest{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: "tx_hash",
		},
	})

	var vout int64 = 0
	assert.Equal(t, &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: "tx_hash",
		},
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index:        0,
					NetworkIndex: &vout,
				},
				Type: defichain.InputOpType,
			},
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 1,
				},
				Type: defichain.OutputOpType,
			},
		},
	}, memTransaction.Transaction)

	mockClient.AssertExpectations(t)
}
