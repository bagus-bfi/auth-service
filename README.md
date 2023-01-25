Go backend repository template with PostgreSQL as the default database.

## Run Locally

Note: We use [docker](https://github.com/docker/cli) and [docker-compose](https://github.com/docker/compose) to run all required infrastructure objects locally.

After git clone, use following commands sequence:

- `make vendor`. To download Go module dependencies into `vendor` folder. (first time only, or on each code dependency changes).
- `make env`. To copy `.env.example` to `.env`. 
  This is for local reference, and will be `exported` right before run the server (see [docker-compose.yml](./docker-compose.yml)). (first time only, or on each config changes).
- `make server/start`. To start the servers (infrastructure and gRPC server).
  - We can use [lazydocker](https://github.com/jesseduffield/lazydocker) to view actively running docker containers.
- `make server/restart`: To restart gRPC server (also restart db-schema to apply if any new db changes).

### Exposed Ports

- [8080](http://0.0.0.0:8080): gRPC HTTP REST proxy server.
  - [OpenAPI Documentation](http://0.0.0.0:8080/openapi-ui/index.html): OpenAPI documentation.
- [8081](http://0.0.0.0:8081): gRPC server.
- [9000](http://0.0.0.0:9000): Adminer web-based postgres web UI. (host: `postgres`, username: `root`, password: `root`)
- _add more exposed ports here_

To stop:

- `make server/stop`. To stop the servers (infrastructure and gRPC server).

See [Makefile](./Makefile) content on what those command does.

### Configuration

All configuration are loaded from environment variables. 
See [env.example](./.env.example) for required environment variables to be set before executing the app.

### Test and Static Analysis

- To run Go unit test, use `make test`. It will also producing coverage results (`coverage.out` and `coverage.html`).
- To run linter and static analysis, use `make lint`.
  It will execute [golangci-lint](https://github.com/golangci/golangci-lint) and use [.golangci.yml](./.golangci.yml) for the rules.

### Build

- `make build/grpc`: to statically build the gRPC server and create a new `build/*` file (default output is for linux amd64).
- `make docker/db-schema`: to build and create a docker image for [Atlas](https://atlasgo.io/) declarative database migration.
- `make docker/grpc`: to build and create a docker image for gRPC server.

### Development

See [Development Notes](./DEVELOPMENT.md).
