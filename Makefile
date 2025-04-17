.PHONY: create-binaries

create-binaries:
	mkdir -p bin
		# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/git-manager.so ./main.go
		# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/git-manager.dylib ./main.go
		# macOS AMD64 (Intel)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/git-manager-amd64.dylib ./main.go
		# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/git-manager.dll ./main.go
