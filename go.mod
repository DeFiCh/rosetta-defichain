module github.com/DeFiCh/rosetta-defichain

go 1.13

require (
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/coinbase/rosetta-sdk-go v0.6.5
	github.com/dgraph-io/badger/v2 v2.2007.2
	github.com/google/addlicense v0.0.0-20200906110928-a0294312aa76 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/rs/zerolog v1.20.0 // indirect
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.16.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/tools v0.0.0-20200904185747-39188db58858 // indirect
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
)

replace github.com/btcsuite/btcd v0.21.0-beta => github.com/DeFiCh/dfid v0.21.0-beta
