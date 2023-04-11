build-content:
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o content cmd/content/main.go
  mv content $(ARTIFACTS_DIR)

build-contentPut:
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o contentPut cmd/contentPut/main.go
  mv contentPut $(ARTIFACTS_DIR)

build-files:
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o files cmd/files/main.go
  mv files $(ARTIFACTS_DIR)

build-filePost:
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o filePost cmd/filePost/main.go
  mv filePost $(ARTIFACTS_DIR)

build-fileDelete:
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o fileDelete cmd/fileDelete/main.go
  mv fileDelete $(ARTIFACTS_DIR)