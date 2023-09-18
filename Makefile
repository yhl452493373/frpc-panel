export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

build: frpc-panel
	cp ./config/frpc-panel.toml ./bin/frpc-panel.toml
	cp -r ./assets/ ./bin/assets/

frpc-panel:
	rm -rf ./bin
	go build -o ./bin/frpc-panel ./cmd/frpc-panel
