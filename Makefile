.PHONY: build doc fmt lint run test vet collect-cover-data test-cover-html test-cover-func
export GO15VENDOREXPERIMENT=1
default: build
build: fmt 
	go build -v -o omega-billing ./
doc:
	godoc -http=:6060 -index
run: build
	./omega-billing
fmt:
	go fmt ./...
test:
	go test ./...


PACKAGES = $(shell find ./ -type d -not -path '*/\.*')
collect-cover-data:
	echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES),\
		go test -v -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)

test-cover-html:
	go tool cover -html=coverage-all.out -o coverage.html

test-cover-func:
	go tool cover -func=coverage-all.out
