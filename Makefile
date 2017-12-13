travis:
	go get -t ./...
	go test ./... --cover -v

test:
	go test ./... --cover -v