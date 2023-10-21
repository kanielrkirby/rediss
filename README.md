# Rediss

Rediss is a simple Redis clone, written completely in Go.

## Features

- [x] TCP server, supports multiple concurrent clients.
- [x] RESP protocol for client-server communication.
- [x] AOF persistence using RESP.
- [ ] Publish/Subscribe (it exists, but has issues with multiple clients, which kind of defeats the purpose).
- [x] Client CLI built with Cobra to communicate with the server.
- [x] Server CLI built with Cobra to handle the database, AOF, and requests from clients.
- [x] Makefile to handle building, as well as the dev servers.

## Usage

### Building

#### To build the rediss-server and rediss-cli binaries

```bash
$ make build
```

#### To run the rediss-server and rediss-cli binaries

```bash
$ rediss-cli
$ rediss-server
```

#### Supported flags for the rediss-server

```bash
$ rediss-server --help
$ redis-server --host localhost
$ rediss-server --port 8080
$ rediss-server --aof /path/to/aof/file
```

#### Supported flags for the rediss-cli

```bash
$ rediss-cli --help
$ rediss-cli --host localhost
$ rediss-cli --port 8080
$ rediss-cli --command GET key
```

### Running the dev servers

#### To run the dev server for the server-side

```bash
$ make dev-server
```

#### To run the dev server for the client-side

```bash
$ make dev-client
```

## References and Honorable Mentions

There were many great resources that helped me build this project. Here are some of them:

- [Redis Docs](https://redis.io/docs/)
- [Go Redis Github](https://github.com/redis/go-redis/)
- [Build Redis from Scratch](https://github.com/ahmedash95/build-redis-from-scratch/)
- [Redis Github](https://github.com/redis/redis/)
- And of course, [Sina](https://sina.khodaveisi.com/), who challenged me to a competition to build a Redis clone better than his in a week. Here's (his version in C++)[https://github.com/sinasun/redis-from-scratch-cpp].
