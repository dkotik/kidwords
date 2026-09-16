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
	@output=$$(go test ./service/http -livePort=8080) || echo "$$output"
generate:
	@clear
	@output=$$(go generate ./...) || echo "$$output"
	@#output=$$(go generate ./... && go test . -update) || echo "$$output"
	@date +"[ %T ]"
build:
	cd cmd/kidwords && goreleaser release --snapshot --clean
install:
	cd ./cmd/kidwords && go build -trimpath -o ~/.local/bin/kidwords
	chmod +x ~/.local/bin/kidwords
