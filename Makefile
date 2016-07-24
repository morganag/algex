export GOPATH=$(shell /usr/bin/pwd)

tests:
	go test algex/factor
