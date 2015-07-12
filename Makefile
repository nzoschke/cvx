install:
	go get ./cli/convox

test:
	go test -v ./...

vendor:
	godep save -r ./...
