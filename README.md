# Static page to IPFS gateway

[![Build Status](https://travis-ci.org/romainbou/scrape-to-ipfs.svg?branch=master)](https://travis-ci.org/romainbou/scrape-to-ipfs)
[![Go Report Card](https://goreportcard.com/badge/github.com/romainbou/scrape-to-ipfs)](https://goreportcard.com/report/github.com/romainbou/scrape-to-ipfs)

Scapes a static web pages and automatically push and serve them via IPFS

## Build and run with Docker

 ```bash
 make image
 make run-docker
 ```

## Test

 ```bash
 make test
 ```

## Environement variables

|  Env        | Default                       |  Description          |
|-------------|-------------------------------|-----------------------|
| PORT        | 8000                          | Server listening port |
| GATEWAY_URL | https://gateway.ipfs.io/ipfs/ | IPFS HTTPS Gateway    |
| IPFS_DAEMON | localhost:5001                | IPFS daemon address   |
| TEMPLATES   | templates/                    | Templates folder      |

## TODO

- [ ] Scrape a page on request and replace or save all relative URLs
- [ ] Push the files to an IPFS daemon
- [x] Redirect to an IPFS gateway
- [ ] Publish the result to IPNS