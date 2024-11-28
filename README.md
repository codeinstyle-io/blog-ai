# Captain

Captain is an experimental blog engine entirely written by AI models. This project explores the possibilities and limitations of AI-generated software by creating a fully functional blogging platform. Every line of code, from the architecture decisions to the implementation details, has been crafted through AI assistance.

## Project Structure

Captain is written in Go and uses:
- Sqlite for data storage
- Gin framework for HTTP routing
- Templ for HTML templates

## Available Make Commands

### Build & Run
- `make build` - Builds the binary in dist/captain
- `make run` - Builds and runs the server
- `make dev` - Runs the server in development mode with live reload
- `make clean` - Removes build artifacts

### Docker
- `make docker-build` - Builds Docker image tagged as captain:latest
- `make docker-run` - Runs the container and exposes port 8080

To use Docker:

Then open http://localhost:8080 in your browser

### Testing & Quality
- `make test` - Runs all unit tests
- `make test-coverage` - Generates test coverage report in coverage.html
- `make lint` - Runs golangci-lint checks
- `make fmt` - Formats Go code

### User Management
- `make create-user` - Creates a new user interactively
- `make update-password` - Updates user password

## Configuration

Captain can be configured using either a YAML config file or environment variables.

### Config File (config.yaml)

The config file can be specified using the -c flag:

```sh
captain run -c /path/to/config.yaml
```

Here are the available options:

| Configuration Key       | Description                         | Default Value   |
|-------------------------|-------------------------------------|-----------------|
| `server.host`           | Host address to bind the server     | `localhost`     |
| `server.port`           | Port number for the server          | `8080`          |
| `db.path`               | SQLite database file location       | `blog.db`       |
| `db.log_level`          | Database logging verbosity          | `warn`          |
| `site.chroma_style`     | Syntax highlighting theme           | `paraiso-dark`  |
| `site.timezone`         | Default timezone for dates          | `UTC`           |


### Environment variables

You can also use these environement variables:

| Environment Variable          | Description                      | Default Value   | Allowed Values                                                                         |
|-------------------------------|----------------------------------|-----------------|----------------------------------------------------------------------------------------|
| `CAPTAIN_SERVER_HOST`         | Host address to bind the server  | `localhost`     | Any valid hostname or IP                                                               |
| `CAPTAIN_SERVER_PORT`         | Port number for the server       | `8080`          | 1-65535                                                                                |
| `CAPTAIN_DB_PATH`             | SQLite database file location    | `blog.db`       | Any valid file path                                                                    |
| `CAPTAIN_DB_LOG_LEVEL`        | Database logging verbosity       | `warn`          | `silent`, `error`, `warn`, `info`                                                      |
| `CAPTAIN_SITE_CHROMA_STYLE`   | Syntax highlighting theme        | `paraiso-dark`  | [Any Chroma style](https://xyproto.github.io/splash/docs/)                             |
| `CAPTAIN_SITE_TIMEZONE`       | Default timezone for dates       | `UTC`           | [Any TZ database name](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)   |



# Contributing

This project is an AI-first experiment. While all contributions are welcome, we encourage:

1. Using AI assistants (like GitHub Copilot) for code generation
2. Documenting AI-human collaboration in pull requests
3. Sharing insights about AI-assisted development
