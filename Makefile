test:
	go test -v ./...

vendor:
	godep save -r ./...
