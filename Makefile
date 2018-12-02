VERSION         :=      $(shell cat ./VERSION)
PWD 			:=		$(shell pwd)

all: test

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

docker-run:
	docker run -v $(PWD)/templates:/tmpl --env TEMPLATES=/tmpl/ romainbou/scrape-to-ipfs

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

.PHONY: install test fmt release