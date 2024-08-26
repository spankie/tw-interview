# Blockparser

Blockparser is a Go package that implements a parser interface to interact with the Ethereum blockchain. It provides methods to get the current block number, retrieve transactions by address, and subscribe to address updates. Additionally, a REST API is provided to expose these functionalities over HTTP.

## Features

- **Parser Interface**: Implements methods to interact with the Ethereum blockchain.
- **REST API**: Exposes the parser functionalities over HTTP.
- **Makefile**: Provides convenient directives to run, build, and test the code.


## Usage

### Running the Server

First of all make sure you have .env file with the following environment variables set

```bash
export TW_PORT=8000
export TW_LOG_LEVEL=debug
```

To start the server, use the following command:

```sh
make run
```

The server will start and listen on the port set in the .env file or default to `8080`.

API Endpoints:

- `GET` `/block`: Returns the current block number.
- `GET` `/transactions/{address}`: Returns the transactions for the specified address.
- `POST` `/subscribe/{address}`: Subscribes to updates for the specified address.

### Example Requests

Get the current block number:

```bash
curl http://localhost:8080/block
```

Get transactions for an address:

```bash
curl http://localhost:8080/transactions/0x1f9840a85d5af5bf1d1762f925bdaddc4201f984
```

Subscribe to an address:

```bash
curl http://localhost:8080/subscribe/0x1f9840a85d5af5bf1d1762f925bdaddc4201f984
```

### Makefile Commands

run: Runs the server.

```bash
make run
```

build: Builds the project.

```bash
make build
```

test: Runs the tests.

```bash
make test
```