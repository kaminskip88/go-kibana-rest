TEST?=./...
PKG_NAME=kbapi
KIBANA_URL ?= http://127.0.0.1:5601
KIBANA_USERNAME ?= elastic
KIBANA_PASSWORD ?= changeme

all: help


test: fmt
	KIBANA_URL=${KIBANA_URL} KIBANA_USERNAME=${KIBANA_USERNAME} KIBANA_PASSWORD=${KIBANA_PASSWORD} go test $(TEST) -v -count 1 -parallel 1 -race -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) -timeout 120m -run=TestKBAPITestSuite/TestKibanaSaveObjectV2

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./


compose.up: docker-elk
	docker-compose -f docker-elk/docker-compose.yml up -d elasticsearch kibana setup

compose.down: docker-elk
	docker-compose -f docker-elk/docker-compose.yml down -v

compose.logs: docker-elk
	docker-compose -f docker-elk/docker-compose.yml logs -f

docker-elk:
	git clone https://github.com/deviantony/docker-elk.git

.PHONY: test fmt compose.up compose.down compose.logs docker-elk
