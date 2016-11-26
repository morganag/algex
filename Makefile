export GOPATH=$(shell /usr/bin/pwd)

tests:
	go test algex/factor
	go test algex/terms
	go test algex/matrix
