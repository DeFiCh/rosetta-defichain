// Code generated by mockery v1.0.0. DO NOT EDIT.

package indexer

import (
	context "context"

	defichain "github.com/DeFiCh/rosetta-defichain/defichain"

	mock "github.com/stretchr/testify/mock"

	types "github.com/coinbase/rosetta-sdk-go/types"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// GetRawBlock provides a mock function with given fields: _a0, _a1
func (_m *Client) GetRawBlock(_a0 context.Context, _a1 *types.PartialBlockIdentifier) (*defichain.Block, []string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *defichain.Block
	if rf, ok := ret.Get(0).(func(context.Context, *types.PartialBlockIdentifier) *defichain.Block); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*defichain.Block)
		}
	}

	var r1 []string
	if rf, ok := ret.Get(1).(func(context.Context, *types.PartialBlockIdentifier) []string); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *types.PartialBlockIdentifier) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetRawTransaction provides a mock function with given fields: ctx, txid, blockhash
func (_m *Client) GetRawTransaction(ctx context.Context, txid string, blockhash string) (*defichain.Transaction, error) {
	ret := _m.Called(ctx, txid, blockhash)

	var r0 *defichain.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *defichain.Transaction); ok {
		r0 = rf(ctx, txid, blockhash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*defichain.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, txid, blockhash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransaction provides a mock function with given fields: ctx, txid
func (_m *Client) GetTransaction(ctx context.Context, txid string) ([]byte, error) {
	ret := _m.Called(ctx, txid)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, txid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkStatus provides a mock function with given fields: _a0
func (_m *Client) NetworkStatus(_a0 context.Context) (*types.NetworkStatusResponse, error) {
	ret := _m.Called(_a0)

	var r0 *types.NetworkStatusResponse
	if rf, ok := ret.Get(0).(func(context.Context) *types.NetworkStatusResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.NetworkStatusResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseBlock provides a mock function with given fields: _a0, _a1, _a2
func (_m *Client) ParseBlock(_a0 context.Context, _a1 *defichain.Block, _a2 map[string]*types.AccountCoin) (*types.Block, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *types.Block
	if rf, ok := ret.Get(0).(func(context.Context, *defichain.Block, map[string]*types.AccountCoin) *types.Block); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Block)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *defichain.Block, map[string]*types.AccountCoin) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PruneBlockchain provides a mock function with given fields: _a0, _a1
func (_m *Client) PruneBlockchain(_a0 context.Context, _a1 int64) (int64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
