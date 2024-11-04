compile:
	go build -o ./aws-mfa ./cmd/

test:
	go test ./...
