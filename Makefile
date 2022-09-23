build:
	go build -o gpr main.go

test:
	go test ./...

detail:
	go test -v ./...
