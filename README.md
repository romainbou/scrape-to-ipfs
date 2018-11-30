# Static page to IPFS gateway

Scapes a static web pages and automatically push and serve them via IPFS

## Environement variables

|  Env        | Default                       |  Description          |
|-------------|-------------------------------|-----------------------|
| PORT        | 8000                          | Server listening port |
| GATEWAY_URL | https://gateway.ipfs.io/ipfs/ |  IPFS HTTPS Gateway   |
| IPFS_DAEMON | localhost:5001                | IPFS daemon address   |

## TODO

[] Scrape a page on request and replace or save all relative URLs
[] Push the files to an IPFS daemon
[] Redirect to an IPFS gateway