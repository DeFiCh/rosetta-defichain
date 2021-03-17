.PHONY: deps build run lint mocks run-mainnet-online run-mainnet-offline run-testnet-online \
	run-testnet-offline check-comments add-license check-license shorten-lines test \
	coverage spellcheck salus build-local coverage-local format check-format check-data-api \
	check-construction-api check-testnet-rosetta-cli-configuration

ADDLICENSE_CMD=go run github.com/google/addlicense
ADDLICENCE_SCRIPT=${ADDLICENSE_CMD} -c "Coinbase, Inc." -l "apache" -v
SPELLCHECK_CMD=go run github.com/client9/misspell/cmd/misspell
GOLINES_CMD=go run github.com/segmentio/golines
GOLINT_CMD=go run golang.org/x/lint/golint
GOVERALLS_CMD=go run github.com/mattn/goveralls
GOIMPORTS_CMD=go run golang.org/x/tools/cmd/goimports
GO_PACKAGES=./services/... ./indexer/... ./defichain/... ./configuration/...
GO_FOLDERS=$(shell echo ${GO_PACKAGES} | sed -e "s/\.\///g" | sed -e "s/\/\.\.\.//g")
TEST_SCRIPT=go test ${GO_PACKAGES}
LINT_SETTINGS=golint,misspell,gocyclo,gocritic,whitespace,goconst,gocognit,bodyclose,unconvert,lll,unparam
PWD=$(shell pwd)
NOFILE=100000
CONTAINER_NAME = rosetta-defichain


deps:
	go get ./...

build:
	docker build -t rosetta-defichain:latest https://github.com/DeFiCh/rosetta-defichain.git

build-local:
	docker build -t rosetta-defichain:latest .

build-release:
	# make sure to always set version with vX.X.X
	docker build -t rosetta-defichain:$(version) .;
	docker save rosetta-defichain:$(version) | gzip > rosetta-defichain-$(version).tar.gz;

run-mainnet-online:
	docker run -d --ulimit --name ${CONTAINER_NAME} "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/data:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -p 8080:8080 -p 8333:8333 -p 8554:8554 -p 8555:8555 rosetta-defichain:latest

run-mainnet-offline:
	docker run -d --rm --name ${CONTAINER_NAME} -e "MODE=OFFLINE" -e "NETWORK=MAINNET" -e "PORT=8081" -p 8081:8081 rosetta-defichain:latest

run-testnet-online:
	docker run -d --rm --name ${CONTAINER_NAME} --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/data:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -p 8080:8080 -p 18333:18333 rosetta-defichain:latest

run-testnet-offline:
	docker run -d --rm --name ${CONTAINER_NAME} -e "MODE=OFFLINE" -e "NETWORK=TESTNET" -e "PORT=8081" -p 8081:8081 rosetta-defichain:latest

train:
	./zstd-train.sh $(network) transaction $(data-directory)

check-comments:
	${GOLINT_CMD} -set_exit_status ${GO_FOLDERS} .

lint: | check-comments
	golangci-lint run --timeout 2m0s -v -E ${LINT_SETTINGS},gomnd

add-license:
	${ADDLICENCE_SCRIPT} .

check-license:
	${ADDLICENCE_SCRIPT} -check .

shorten-lines:
	${GOLINES_CMD} -w --shorten-comments ${GO_FOLDERS} .

format:
	gofmt -s -w -l .
	${GOIMPORTS_CMD} -w .

check-format:
	! gofmt -s -l . | read
	! ${GOIMPORTS_CMD} -l . | read

test:
	${TEST_SCRIPT}

coverage:
	if [ "${COVERALLS_TOKEN}" ]; then ${TEST_SCRIPT} -coverprofile=c.out -covermode=count; ${GOVERALLS_CMD} -coverprofile=c.out -repotoken ${COVERALLS_TOKEN}; fi

coverage-local:
	${TEST_SCRIPT} -cover

salus:
	docker run --rm -t -v ${PWD}:/home/repo coinbase/salus

spellcheck:
	${SPELLCHECK_CMD} -error .

mocks:
	rm -rf mocks;
	mockery --dir indexer --all --case underscore --outpkg indexer --output mocks/indexer;
	mockery --dir services --all --case underscore --outpkg services --output mocks/services;
	${ADDLICENCE_SCRIPT} .;

# Possible configuration parameters:
# ----------------------------------
# GENERAL:
# - https://pkg.go.dev/github.com/coinbase/rosetta-cli/configuration#Configuration
# 
# DATA:
# - https://pkg.go.dev/github.com/coinbase/rosetta-cli/configuration#DataConfiguration
# 
# CONSTRUCTION:
# - https://pkg.go.dev/github.com/coinbase/rosetta-cli/configuration#ConstructionConfiguration
# 
# NOTE: in order to perform all checks correctly rosetta-cli should be installed from 
# the following source:
# 	- https://github.com/DeFiDev/rosetta-cli/tree/v0.6.7-1
# 
# This is a fork of original rosetta-cli implementation, but with some fix applied. 
# Detailed fix description may be found here:
# 	- https://github.com/coinbase/rosetta-cli/issues/222
#  
check-data-api: check-testnet-rosetta-cli-configuration
	rosetta-cli check:data --configuration-file ${CLI_TESTNET_CONFIG_PATH}

check-construction-api: check-testnet-rosetta-cli-configuration
	rosetta-cli check:construction --configuration-file ${CLI_TESTNET_CONFIG_PATH}

# NOTE: may be provided as an environment variable
CLI_TESTNET_CONFIG_PATH = `pwd`/rosetta-cli-conf/testnet/config.json

check-testnet-rosetta-cli-configuration:
	rosetta-cli configuration:validate ${CLI_TESTNET_CONFIG_PATH}