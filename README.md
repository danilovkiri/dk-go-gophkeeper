# Gophkeeper

Gophkeeper is a data storage system written in Go and implementing:
1. TUI client
2. Secure server data storage
3. GRPC for data transfer
4. Login/Password authorization

## Usage

### Tuning

All parameters can be passed via environment variables:
1. SERVER_ADDRESS — a string representation of a server address
2. DATABASE_DSN — a string representation of a PSQL DB DSN
3. USER_KEY — a special key for ciphering data on the server side (default `jds__63h3_7ds`)
4. BEARER_KEY — a GRPC context metadata key to be used in authorization (default `token`)
5. HANDLERS_TO — a shared timeout for server unary operations (in ms, default `500`)

### Server

Compile both [client](./cmd/client/main.go) and [server](./cmd/server/main.go) using the same configuration parameters
for compatibility.

Run the server application (or compiled binary) supplying it with command line arguments:

```shell
go run ./cmd/server/main.go -a :8080 -c ./config.json -d postgres://user:password@localhost:5432/gophkeeper
```

where `-a` defines a GRPC server address, `-d` — PSQL DB DSN, `-c` — a path to configuration file. The example of
configuration file is shown below:

```json
{
  "server_address": ":8080",
  "database_dsn": "postgres://user:password@localhost:5432/gophkeeper"
}
```

Note that command-line arguments are prioritized over environment and JSON config-derived arguments. Environment
variables have, in turn, higher priority that JSON config-derived ones.

### Client

Run the TUI application (or compiled binary):

```shell
go run ./cmd/client/main.go
```

The resizable window will appear inside the console:

<img src="./resources/mainView.png" alt="drawing" width="700"/>

### Notes

1. A user must log in or register first, no data can be stored unless authentication completed
2. Identifiers in the data addition forms must be unique for each type of data
3. Any data addition/removal immediately sends requests to the server
4. An existing user must press the `Sync` button after logging in and prior to data addition; syncing retrieves all user
data from the server and overwrites any local changes.
5. Any errors will be reported in the bottom part of the screen.





