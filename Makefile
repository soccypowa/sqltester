test:
	go test -v -cover -short ./...

tidy:
	go mod tidy

build:
	GOOS='windows' GOARCH='amd64' go build .

clean:
	rm -rf sqltester.exe
	rm -rf sqltester

listdepupdates:
	go list -m -u all

updatedeps:
	go get -u

.PHONY: test tidy build listdepupdates updatedeps