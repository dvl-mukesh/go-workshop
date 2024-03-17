test:
	@echo "[+] Running tests"
	@go test -v ./... -count=1 -race -covermode=atomic -coverprofile=coverage.out 

cover: test
	@echo "[+] Generating coverage report"
	@go tool cover -html=coverage.out

.PHONY: test cover