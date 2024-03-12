GOBIN ?= $$(go env GOPATH)/bin

install:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/tools/gopls@latest

doc:
	rm -rf docs
	swag fmt
	swag init --dir cmd/server/,api,internals

build:
	go build -o ./bin/server ./cmd/server/main.go

dev:
	./bin/air -d

run:
	go run ./cmd/server/main.go

get:
	go get ./cmd/server

dbuild:
	docker-compose build api

dup:
	docker-compose up -d api

dlogs:
	docker-compose logs --tail=50 -f api

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml

show-coverage:
	go tool cover -html=cover.out