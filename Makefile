run-server:
	go run main.go
.PHONY: run-server

run-opa:
	opa run --server --addr localhost:8081 --watch policies
.PHONY: run-opa
