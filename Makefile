GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

all: clean fmt deps build move

build: mac windows linux

dev: clean fmt deps mac move

move:
	tar -xvf  bin/packer-builder-alicloud-ecs_darwin-amd64.tgz && mv bin/packer-builder-alicloud-ecs   $(shell dirname `which packer`)

test: 
	PACKER_ACC=1 go test -v ./ecs -timeout 120m

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


mac: deps
	GOOS=darwin GOARCH=amd64 go build -o bin/packer-builder-alicloud-ecs
	GOOS=darwin GOARCH=amd64 go build -o bin/packer-post-processor-alicloud-import
	tar czvf bin/packer-builder-alicloud-ecs_darwin-amd64.tgz bin/packer-builder-alicloud-ecs bin/packer-post-processor-alicloud-import
	rm -rf bin/packer-provider-alicloud-ecs bin/packer-post-processor-alicloud-import

windows: deps
	GOOS=windows GOARCH=amd64 go build -o bin/packer-builder-alicloud-ecs.exe
	GOOS=windows GOARCH=amd64 go build -o bin/packer-post-processor-alicloud-import.exe
	tar czvf bin/packer-builder-alicloud-ecs_windows-amd64.tgz bin/packer-builder-alicloud-ecs.exe bin/packer-post-processor-alicloud-import.exe
	rm -rf bin/packer-builder-alicloud-ecs.exe bin/packer-post-processor-alicloud-import.exe

linux: deps
	GOOS=linux GOARCH=amd64 go build -o bin/packer-builder-alicloud-ecs
	GOOS=linux GOARCH=amd64 go build -o bin/packer-post-processor-alicloud-import
	tar czvf bin/packer-builder-alicloud-ecs_linux-amd64.tgz bin/packer-builder-alicloud-ecs bin/packer-post-processor-alicloud-import
	rm -rf bin/packer-builder-alicloud-ecs bin/packer-post-processor-alicloud-import

clean:
	rm -rf bin/*
