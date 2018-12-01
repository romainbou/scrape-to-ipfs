VERSION         :=      $(shell cat ./VERSION)

all: install

install:
	go install -v

run:
	go run server.go

test:
	go test ./... -v

fmt:
	go fmt ./... -v

image:
	docker build -t romainbou/scrape-to-ipfs .

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

.PHONY: install test fmt release