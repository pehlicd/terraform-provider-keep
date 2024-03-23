TEST?=$$(go list ./... | grep -v 'vendor')
BINARY=terraform-provider-keep

default: install

build:
	go build -o ${BINARY}

install: build
	mv ${BINARY} /Users/pehlivan/.terraform.d/plugins/terraform.local/local/keep/1.0.0/darwin_arm64/${BINARY}_v1.0.0

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
	xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m