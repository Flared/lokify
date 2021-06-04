.PHONY: go-build
go-build:
	go build -o build/lokify cmd/api/main.go

.PHONY: go-run
go-run:
	go run cmd/api/main.go

.PHONY: go-tests
go-tests:
	go test ./pkg/...

ui/node_modules: ui/package.json \
                 ui/package-lock.json
	rm -rf ui/node_modules
	npm --prefix ui ci

.PHONY: clean 
clean:
	rm -rf ui/node_modules

.PHONY: npm-start
npm-start: ui/node_modules
	npm --prefix ui start
