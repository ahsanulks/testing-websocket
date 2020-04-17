dep:
	go mod download

run: dep
	go run main/main.go