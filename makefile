all: feishusign

feishusign: ./cmd/main.go
	rm -rf ./logs
	go build -o feishusign ./cmd/main.go 
