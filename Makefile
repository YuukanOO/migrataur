install:
	go get -t ./...

test:
	go test ./... --cover -v