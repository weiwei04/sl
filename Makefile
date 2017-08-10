test:
	go test -cover $(shell go list ./... | grep -v /vendor/)
