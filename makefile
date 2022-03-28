build-amd:
	GOOS=linux GOARCH=amd64 go build -o bin/bootstrap main.go
	cd bin && zip function.zip bootstrap

build-arm:
	GOOS=linux GOARCH=arm64 go build -o bin/bootstrap main.go
	cd bin && zip function.zip bootstrap