install:
	go get ./...

server:
	rerun -build github.com/nzoschke/cvx

test:
	go test -v ./...

vendor:
	godep save -r ./...
