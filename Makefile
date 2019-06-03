generate: Makefile
	docker build -t alerting -f ./Dockerfile .

dev:
	rm -f alert
	echo "Building binary..."
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -ldflags '-w' -o ./alert cmd/alert/main.go
	echo "Building images..."
	docker build -t alerting -f ./Dockerfile.dev .
	echo "Built successfully"

clean:
	rm -f alert