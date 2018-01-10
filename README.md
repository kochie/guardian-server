<img src='https://d.pr/SLO9Sg/3123P4IP+' />

# guardian-server [![Build Status](https://travis-ci.org/kochie/guardian-server.svg?branch=master)](https://travis-ci.org/kochie/guardian-server) [![Coverage Status](https://coveralls.io/repos/github/kochie/guardian-server/badge.svg?branch=master)](https://coveralls.io/github/kochie/guardian-server?branch=master)
Server for communicating with guardian hosts.

# Building and Installation
The project can be build on either docker or bare metal, however most development is done with the provided dockerfile.

## Bare Metal

- Clone the repository.

```bash
git clone git@github.com:kochie/guardian-server.git
```
- Start a [redis](https:/redis.io) instance somewhere and take note of the connection settings.
- Start a [postgres](https://www.postgresql.org/) database somewhere and take note of the connection settings as well.
- Modify the environment variables found in `config.json` and update any settings required.
- Build the service with go, normally a sequence like.

```bash
go test -v ./...
go build main.go
```

- [Optionally] Install the binary.
```bash
go install main.go
```

## Docker
To build in docker compose.
```bash
docker-compose build 
```
To run in docker compose.
```bash
docker-compose up
```
