build-content:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o content main.go
	mv content $(ARTIFACTS_DIR)

build-article:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o article main.go
	mv article $(ARTIFACTS_DIR)

build-contentPut:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o contentPut main.go
	mv contentPut $(ARTIFACTS_DIR)

build-files:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o files main.go
	mv files $(ARTIFACTS_DIR)

build-filePost:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o filePost main.go
	mv filePost $(ARTIFACTS_DIR)

build-fileDelete:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o fileDelete main.go
	mv fileDelete $(ARTIFACTS_DIR)