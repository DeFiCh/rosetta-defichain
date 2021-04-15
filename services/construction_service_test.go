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
	"encoding/hex"
	"testing"

	"github.com/DeFiCh/rosetta-defichain/configuration"
	"github.com/DeFiCh/rosetta-defichain/defichain"
	mocks "github.com/DeFiCh/rosetta-defichain/mocks/services"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func forceHexDecode(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("could not decode hex %s", s)
	}

	return b
}

func forceMarshalMap(t *testing.T, i interface{}) map[string]interface{} {
	m, err := types.MarshalMap(i)
	if err != nil {
		t.Fatalf("could not marshal map %s", types.PrintStruct(i))
	}

	return m
}

func TestConstructionService(t *testing.T) {
	networkIdentifier = &types.NetworkIdentifier{
		Network:    defichain.TestnetNetwork,
		Blockchain: defichain.Blockchain,
	}

	cfg := &configuration.Configuration{
		Mode:     configuration.Online,
		Network:  networkIdentifier,
		Params:   defichain.TestnetParams,
		Currency: defichain.TestnetCurrency,
	}

	mockIndexer := &mocks.Indexer{}
	mockClient := &mocks.Client{}
	servicer := NewConstructionAPIService(cfg, mockClient, mockIndexer)
	ctx := context.Background()

	// Test Derive
	publicKey := &types.PublicKey{
		Bytes: forceHexDecode(
			t,
			"03d3d13e8180b10dfed0db4db9fb9013a5dcbdab64dff45d16310313c2929e71ac",
		),
		CurveType: types.Secp256k1,
	}
	deriveResponse, err := servicer.ConstructionDerive(ctx, &types.ConstructionDeriveRequest{
		NetworkIdentifier: networkIdentifier,
		PublicKey:         publicKey,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{
			Address: "tf1q4mfg3kqm4c7f8j7tyy028fa4uaw7zdzyerq5qx",
		},
	}, deriveResponse)

	// Test Preprocess
	ops := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: defichain.InputOpType,
			Account: &types.AccountIdentifier{
				Address: "tf1q4mfg3kqm4c7f8j7tyy028fa4uaw7zdzyerq5qx",
			},
			Amount: &types.Amount{
				Value:    "-999999444",
				Currency: defichain.TestnetCurrency,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "d435e59999faffb36890651bd0cfd9795bb06bfe96d5ce62f1acd9cc1ce12a37:0",
				},
				CoinAction: types.CoinSpent,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			Type: defichain.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "tf1qr4pelt25fk869v2sfrwxp7ttm7jyx0pyn4dxph",
			},
			Amount: &types.Amount{
				Value:    "954843",
				Currency: defichain.TestnetCurrency,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 2,
			},
			Type: defichain.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "tf1qmzzp8v6q6pxanuvaey49c6pkcmrvdxudun0v8j",
			},
			Amount: &types.Amount{
				Value:    "44657",
				Currency: defichain.TestnetCurrency,
			},
		},
	}
	feeMultiplier := float64(0.75)
	preprocessResponse, err := servicer.ConstructionPreprocess(
		ctx,
		&types.ConstructionPreprocessRequest{
			NetworkIdentifier:      networkIdentifier,
			Operations:             ops,
			SuggestedFeeMultiplier: &feeMultiplier,
		},
	)
	assert.Nil(t, err)
	options := &preprocessOptions{
		Coins: []*types.Coin{
			{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "d435e59999faffb36890651bd0cfd9795bb06bfe96d5ce62f1acd9cc1ce12a37:0",
				},
				Amount: &types.Amount{
					Value:    "-999999444",
					Currency: defichain.TestnetCurrency,
				},
			},
		},
		EstimatedSize: 142,
		FeeMultiplier: &feeMultiplier,
	}
	assert.Equal(t, &types.ConstructionPreprocessResponse{
		Options: forceMarshalMap(t, options),
	}, preprocessResponse)

	// Test Metadata
	metadata := &constructionMetadata{
		ScriptPubKeys: []*defichain.ScriptPubKey{
			{
				ASM:          "0 aed288d81bae3c93cbcb211ea3a7b5e75de13444",
				Hex:          "0014aed288d81bae3c93cbcb211ea3a7b5e75de13444",
				RequiredSigs: 1,
				Type:         "witness_v0_keyhash",
				Addresses: []string{
					"tf1q4mfg3kqm4c7f8j7tyy028fa4uaw7zdzyerq5qx",
				},
			},
		},
	}

	// Normal Fee
	mockIndexer.On(
		"GetScriptPubKeys",
		ctx,
		options.Coins,
	).Return(
		metadata.ScriptPubKeys,
		nil,
	).Once()
	mockClient.On(
		"SuggestedFeeRate",
		ctx,
		defaultConfirmationTarget,
	).Return(
		defichain.MinFeeRate*10,
		nil,
	).Once()
	metadataResponse, err := servicer.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
		NetworkIdentifier: networkIdentifier,
		Options:           forceMarshalMap(t, options),
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionMetadataResponse{
		Metadata: forceMarshalMap(t, metadata),
		SuggestedFee: []*types.Amount{
			{
				Value:    "1065", // 1,420 * 0.75
				Currency: defichain.TestnetCurrency,
			},
		},
	}, metadataResponse)

	// Low Fee
	mockIndexer.On(
		"GetScriptPubKeys",
		ctx,
		options.Coins,
	).Return(
		metadata.ScriptPubKeys,
		nil,
	).Once()
	mockClient.On(
		"SuggestedFeeRate",
		ctx,
		defaultConfirmationTarget,
	).Return(
		defichain.MinFeeRate,
		nil,
	).Once()
	metadataResponse, err = servicer.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
		NetworkIdentifier: networkIdentifier,
		Options:           forceMarshalMap(t, options),
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionMetadataResponse{
		Metadata: forceMarshalMap(t, metadata),
		SuggestedFee: []*types.Amount{
			{
				Value:    "142", // we don't go below minimum fee rate
				Currency: defichain.TestnetCurrency,
			},
		},
	}, metadataResponse)

	// Test Payloads
	unsignedRaw := "7b227472616e73616374696f6e223a2230313030303030303031333732616531316363636439616366313632636564353936666536626230356237396439636664303162363539303638623366666661393939396535333564343030303030303030303066666666666666663032646239313065303030303030303030303136303031343164343339666164353434643866613262313530343864633630663936626466613434333363323437316165303030303030303030303030313630303134643838343133623334306430346464396631396463393261356336383336633663366336396238643030303030303030222c227363726970745075624b657973223a5b7b2261736d223a22302061656432383864383162616533633933636263623231316561336137623565373564653133343434222c22686578223a223030313461656432383864383162616533633933636263623231316561336137623565373564653133343434222c2272657153696773223a312c2274797065223a227769746e6573735f76305f6b657968617368222c22616464726573736573223a5b2274663171346d6667336b716d34633766386a37747979303238666134756177377a647a79657271357178225d7d5d2c22696e7075745f616d6f756e7473223a5b222d393939393939343434225d2c22696e7075745f616464726573736573223a5b2274663171346d6667336b716d34633766386a37747979303238666134756177377a647a79657271357178225d7d" // nolint
	payloadsResponse, err := servicer.ConstructionPayloads(ctx, &types.ConstructionPayloadsRequest{
		NetworkIdentifier: networkIdentifier,
		Operations:        ops,
		Metadata:          forceMarshalMap(t, metadata),
	})
	val0 := int64(0)
	val1 := int64(1)
	parseOps := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        0,
				NetworkIndex: &val0,
			},
			Type: defichain.InputOpType,
			Account: &types.AccountIdentifier{
				Address: "tf1q4mfg3kqm4c7f8j7tyy028fa4uaw7zdzyerq5qx",
			},
			Amount: &types.Amount{
				Value:    "-999999444",
				Currency: defichain.TestnetCurrency,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "d435e59999faffb36890651bd0cfd9795bb06bfe96d5ce62f1acd9cc1ce12a37:0",
				},
				CoinAction: types.CoinSpent,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        1,
				NetworkIndex: &val0,
			},
			Type: defichain.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "tf1qr4pelt25fk869v2sfrwxp7ttm7jyx0pyn4dxph",
			},
			Amount: &types.Amount{
				Value:    "954843",
				Currency: defichain.TestnetCurrency,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        2,
				NetworkIndex: &val1,
			},
			Type: defichain.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "tf1qmzzp8v6q6pxanuvaey49c6pkcmrvdxudun0v8j",
			},
			Amount: &types.Amount{
				Value:    "44657",
				Currency: defichain.TestnetCurrency,
			},
		},
	}

	assert.Nil(t, err)

	// NOTE: starting from now some test data may be not accurate. But it should be
	// fine for test purposes
	signingPayload := &types.SigningPayload{
		Bytes: forceHexDecode(
			t,
			"69b9d24fb4ee7e69b02ad457c50a56460da918bdbc83e50fbbe69c23164edbc0",
		),
		AccountIdentifier: &types.AccountIdentifier{
			Address: "tf1q4mfg3kqm4c7f8j7tyy028fa4uaw7zdzyerq5qx",
		},
		SignatureType: types.Ecdsa,
	}
	assert.Equal(t, &types.ConstructionPayloadsResponse{
		UnsignedTransaction: unsignedRaw,
		Payloads:            []*types.SigningPayload{signingPayload},
	}, payloadsResponse)

	// Test Parse Unsigned
	parseUnsignedResponse, err := servicer.ConstructionParse(ctx, &types.ConstructionParseRequest{
		NetworkIdentifier: networkIdentifier,
		Signed:            false,
		Transaction:       unsignedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionParseResponse{
		Operations:               parseOps,
		AccountIdentifierSigners: []*types.AccountIdentifier{},
	}, parseUnsignedResponse)

	// Test Combine
	signedRaw := "7b227472616e73616374696f6e223a22303130303030303030303031303133373261653131636363643961636631363263656435393666653662623035623739643963666430316236353930363862336666666139393939653533356434303030303030303030306666666666666666303264623931306530303030303030303030313630303134316434333966616435343464386661326231353034386463363066393662646661343433336332343731616530303030303030303030303031363030313464383834313362333430643034646439663139646339326135633638333663366336633639623864303234373330343430323230323538373665633862396635316433343361356135366163353439633063383238303035656634356562653964613136366462363435633039313537323233663032323034636430386237323738613838383961383131333539313562636531306431656633626239326232313766383161306465376537396666623364666436616335303132313033643364313365383138306231306466656430646234646239666239303133613564636264616236346466663435643136333130333133633239323965373161633030303030303030222c22696e7075745f616d6f756e7473223a5b222d393939393939343434225d7d" // nolint
	combineResponse, err := servicer.ConstructionCombine(ctx, &types.ConstructionCombineRequest{
		NetworkIdentifier:   networkIdentifier,
		UnsignedTransaction: unsignedRaw,
		Signatures: []*types.Signature{
			{
				Bytes: forceHexDecode(
					t,
					"25876ec8b9f51d343a5a56ac549c0c828005ef45ebe9da166db645c09157223f4cd08b7278a8889a81135915bce10d1ef3bb92b217f81a0de7e79ffb3dfd6ac5", // nolint
				),
				SigningPayload: signingPayload,
				PublicKey:      publicKey,
				SignatureType:  types.Ecdsa,
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionCombineResponse{
		SignedTransaction: signedRaw,
	}, combineResponse)

	// Test Parse Signed
	parseSignedResponse, err := servicer.ConstructionParse(ctx, &types.ConstructionParseRequest{
		NetworkIdentifier: networkIdentifier,
		Signed:            true,
		Transaction:       signedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionParseResponse{
		Operations: parseOps,
		AccountIdentifierSigners: []*types.AccountIdentifier{
			{Address: "tf1q4mfg3kqm4c7f8j7tyy028fa4uaw7zdzyerq5qx"},
		},
	}, parseSignedResponse)

	// Test Hash
	transactionIdentifier := &types.TransactionIdentifier{
		Hash: "d72427f967dbe9691328c03cb0c4be74d33dd22b537cd4a081d6ba4da970e8fc",
	}
	hashResponse, err := servicer.ConstructionHash(ctx, &types.ConstructionHashRequest{
		NetworkIdentifier: networkIdentifier,
		SignedTransaction: signedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.TransactionIdentifierResponse{
		TransactionIdentifier: transactionIdentifier,
	}, hashResponse)

	// Test Submit
	deFichainTransaction := "01000000000101372ae11cccd9acf162ced596fe6bb05b79d9cfd01b659068b3fffa9999e535d40000000000ffffffff02db910e00000000001600141d439fad544d8fa2b15048dc60f96bdfa4433c2471ae000000000000160014d88413b340d04dd9f19dc92a5c6836c6c6c69b8d02473044022025876ec8b9f51d343a5a56ac549c0c828005ef45ebe9da166db645c09157223f02204cd08b7278a8889a81135915bce10d1ef3bb92b217f81a0de7e79ffb3dfd6ac5012103d3d13e8180b10dfed0db4db9fb9013a5dcbdab64dff45d16310313c2929e71ac00000000" // nolint
	mockClient.On(
		"SendRawTransaction",
		ctx,
		deFichainTransaction,
	).Return(
		transactionIdentifier.Hash,
		nil,
	)
	submitResponse, err := servicer.ConstructionSubmit(ctx, &types.ConstructionSubmitRequest{
		NetworkIdentifier: networkIdentifier,
		SignedTransaction: signedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.TransactionIdentifierResponse{
		TransactionIdentifier: transactionIdentifier,
	}, submitResponse)

	mockClient.AssertExpectations(t)
	mockIndexer.AssertExpectations(t)
}
