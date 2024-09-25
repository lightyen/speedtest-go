GO_BUILD := go build

LDFLAGS := -s -w

all: linux

linux:
	GOOS=linux $(GO_BUILD) -ldflags="${LDFLAGS}" -o app cmd/main.go

arm7:
	GOOS=linux GOARCH=arm GOARM=7 $(GO_BUILD) -ldflags="${LDFLAGS}" -o app cmd/main.go

arm64:
	GOOS=linux GOARCH=arm64 $(GO_BUILD) -ldflags="${LDFLAGS}" -o app cmd/main.go

windows:
	GOOS=windows GOARCH=amd64 $(GO_BUILD) -ldflags="${LDFLAGS}" -o app.exe cmd/main.go

darwin:
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -ldflags="${LDFLAGS}" -o app cmd/main.go
