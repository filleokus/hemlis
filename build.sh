version="v0.1.0"
commit=$(git rev-parse --short HEAD)
buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

go test ./...
go build -v -ldflags="-X 'main.Version=$version' -X 'main.Commit=$commit' -X 'main.BuildDate=$buildDate'" ./cmd/hemlis-combine