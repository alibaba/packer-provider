GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

all: fmt deps build move

build:
	go build -o packer-builder-alicloud

move:
	mv packer-builder-alicloud   $(shell dirname `which packer`)

test: 
	PACKER_ACC=1 go test -v ./alicloud -timeout 120m

vet:
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"


deps:
	go get -u github.com/kardianos/govendor
	govendor sync
	go get golang.org/x/crypto/curve25519
	go get golang.org/x/crypto/ed25519
