all:
	go get -d -v ./...
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../bin/main .
	chmod +x ../bin/main
