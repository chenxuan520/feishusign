all: feishusign

feishusign: ./cmd/main.go
	go build -o feishusign ./cmd/main.go 
