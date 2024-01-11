test:
	go test -v -cover -short ./...

tidy:
	go mod tidy

build:
	GOOS='windows' GOARCH='amd64' go build .

clean:
	rm -rf sqltester.exe
	rm -rf sqltester

.PHONY: test tidy build