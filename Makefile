-include .env
export

default:
	@clear
	@output=$$(go test ./...) || echo "$$output"
	@date +"[ %T ]"
short:
	@clear
	@output=$$(go test -short ./...) || echo "$$output" | grep -Ev "^(ok|\\?)"
	@date +"[ %T ]"
generate:
	@clear
	@output=$$(go generate ./... && go test . -update) || echo "$$output"
	@date +"[ %T ]"
build:
	cd cmd/kidwords && goreleaser release --snapshot --clean
install:
	cd ./cmd/kidwords && go build -trimpath -o ~/.local/bin/kidwords
	chmod +x ~/.local/bin/kidwords
