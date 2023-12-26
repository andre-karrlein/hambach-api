build-content:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/content/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-article:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/article/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-contentPut:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/contentPut/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-files:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/files/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-filePost:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/filePost/main.go
	mv bootstrap $(ARTIFACTS_DIR)

build-fileDelete:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap cmd/fileDelete/main.go
	mv bootstrap $(ARTIFACTS_DIR)
