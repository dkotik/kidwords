-include .env
export

default:
	@clear
	@output=$$(go test ./...) || echo "$$output"
	@date +"[ %T ]"
short:
	@clear
	@#output=$$(go test -short -run=^TestEncode$$) || echo "$$output" | grep -Ev "^(ok|\\?)"
	@output=$$(go test -short ./...) || echo "$$output" | grep -Ev "^(ok|\\?)"
	@date +"[ %T ]"
live:
	@clear
	@date +"[ %T at http://localhost:8080/ ]"
	@output=$$(cd ./service/http && go generate . && go test . -livePort=8080) || echo "$$output"
generate:
	@clear
	@output=$$(go generate ./...) || echo "$$output"
	@#output=$$(go generate ./... && go test . -update) || echo "$$output"
	@cp service/http/media/wasm.html docs/index.html
	@cp service/http/media/wasm_exec.js docs/
	@cp service/http/media/wkdw.v0.wasm docs/
	@cp service/http/media/bulma.min.css docs/
	@date +"[ %T ]"
build: generate
	cd cmd/kidwords && goreleaser release --snapshot --clean
install: generate default
	cd ./cmd/kidwords && go build -trimpath -o ~/.local/bin/kidwords
	chmod +x ~/.local/bin/kidwords
