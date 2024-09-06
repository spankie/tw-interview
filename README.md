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

# Technical Details

## Design

The project is structured into the following packages:

- **cmd**: Contains the main package to start the web api server.
- **blockparser**: Contains the core logic of the parser.

### Blockparser

The `blockparser` package implements the `Parser` interface, which provides methods to interact
with the Ethereum blockchain using the cloudflare eth api. The cloudflare eth api is abstracted in
an sdk that implements the `BlockchainQuerier` interface used by the `blockparser`. The `blockparser`
also uses a datastore to store the subscribed addresses and their transactions.

The source files for blockparser package are in the same folder because they are really small and related
to each other. The `blockparser.go` file contains the implementation of the `Parser` interface and the
`blockparser_test.go` file contains the tests for the parser. The `db.go` file contains the implementation
of the default datastore used by the parser and the `db_test.go` file contains the tests for the datastore.

The blockparser works by polling the cloudflare eth api at regular intervals to get the latest block
number and filtering transaction in each block from the last scanned block to the latest block.
The transactions that match any of the subscribed addresses are stored in the datastore identified
by the address. The interval for polling is fully configurable and defaults to 1 minute if not set by the user.

Subscribing an address is done by adding the address to the datastore so when the polling runs, transactions
can be checked against the subscribed addresses.

The parser is configurable using options that can be set when creating a new parser instance. The options
include the polling interval, the blockchain querier, and the datastore. This enables the parser to use
different implementation of the blockchain querier and datastore if needed and also enables better testing of
the parser by using mock implementations of the blockchain querier and datastore. The default for these
options is set to use the cloudflare eth api sdk and the memoryStore respectively.

### Datastore

An implementation of the datastore is provided in the `db.go`. The `memoryStore` is a simple in-memory
generic datastore that stores the transactions for each address in a map. The datastore implements the `Datastore`
interface that provides methods to store and retrieve transactions for an address. It is made to be generic so
it can store different data types if needed. The datastore implementation uses a map and a mutex to ensure
thread safety when accessing and mutating the data.

## REST API

A REST api is exposed to interact with the blockchain parser. The api provides endpoint described above.


