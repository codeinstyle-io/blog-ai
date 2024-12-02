# Captain

[![CI](https://github.com/shinuza/captain/actions/workflows/ci.yml/badge.svg)](https://github.com/shinuza/captain/actions/workflows/ci.yml) [![](https://img.shields.io/github/v/release/shinuza/captain)](https://github.com/shinuza/captain/releases) [![](https://img.shields.io/badge/license-MIT-green)](https://github.com/shinuza/captain/blob/master/LICENSE)

Captain is a no-installation AI written blog engine. Just run the binary and you are ready to start writing. No database setup, no dependencies to install. Just pure blogging.

## Trivia

Captain is an experimental blog engine entirely written by AI models. This project explores the possibilities and limitations of AI-generated software by creating a fully functional blogging platform. Every line of code, from the architecture decisions to the implementation details, has been crafted through AI assistance.

## Project Structure

Captain is written in Go and uses:
- Sqlite for data storage
- Gin framework for HTTP routing
- [html/template](https://pkg.go.dev/html/template) for HTML templates


## Installation

### From Binary Releases

1. Download the latest release for your platform from the [releases page](https://github.com/shinuza/captain/releases)
2. Extract the archive:
   ```sh
   unzip captain-<platform>.zip
   ```
3. Move the binary to a location in your PATH:
   ```sh
   sudo mv captain /usr/local/bin/
   ```
4. Verify the installation:
   ```sh
   captain
   ```

### From Source

1. Ensure you have Go 1.21 or later installed
2. Clone the repository:
   ```sh
   git clone https://github.com/shinuza/captain.git
   cd captain
   ```
3. Build and install:
   ```sh
   make build
   sudo mv dist/captain /usr/local/bin/
   ```

### Using Docker

```sh
docker pull shinuza/captain:latest
docker run -p 8080:8080 shinuza/captain:latest
```

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

If the -c flag is not provided, Captain will look for a config file named config.yaml in the current directory or in /etc/captain/

Here are the available options:

| Configuration Key       | Description                         | Default Value   |
|-------------------------|-------------------------------------|-----------------|
| `server.host`           | Host address to bind the server     | `localhost`     |
| `server.port`           | Port number for the server          | `8080`          |
| `db.path`               | SQLite database file location       | `blog.db`       |
| `db.log_level`          | Database logging verbosity          | `warn`          |
| `site.chroma_style`     | Syntax highlighting theme           | `paraiso-dark`  |
| `site.timezone`         | Default timezone for dates          | `UTC`           |
| `site.theme`            | Theme name to use                   | `""`            |
| `site.title`            | Site title                         | `Captain`       |
| `site.subtitle`         | Site subtitle                      | `A simple blog engine` |
| `site.posts_per_page`   | Number of posts per page           | `3`             |

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
| `CAPTAIN_SITE_THEME`          | Theme name to use                | `""`            | Empty string for default theme, or any theme name in themes directory                  |
| `CAPTAIN_SITE_TITLE`          | Site title                       | `Captain`       | Any string                                                                             |
| `CAPTAIN_SITE_SUBTITLE`       | Site subtitle                    | `A simple blog engine` | Any string                                                                      |
| `CAPTAIN_SITE_POSTS_PER_PAGE` | Number of posts per page         | `3`             | Any positive integer                                                                   |

### Themes

Captain supports customizable themes for the public site. The admin interface maintains a consistent look regardless of the theme selected.

#### Default Theme
When `site.theme` is empty (`""`), Captain uses its embedded default theme.

#### Custom Themes
To use a custom theme:

1. Create a directory in `themes/` with your theme name (e.g., `themes/mytheme/`)
2. Add the required theme files:
   ```
   themes/mytheme/
   ├── templates/
   │   ├── header.tmpl
   │   ├── footer.tmpl
   │   └── ... (add your custom templates here)
   └── static/
       ├── css/
       │   └── main.css
       └── js/
           └── main.js
   ```
3. Set `site.theme: "mytheme"` in your config.yaml or `CAPTAIN_SITE_THEME=mytheme`

Custom themes can override any of the default templates and provide their own static assets.

# Contributing

This project is an AI-first experiment. While all contributions are welcome, we encourage:

1. Using AI assistants (like GitHub Copilot) for code generation
2. Documenting AI-human collaboration in pull requests
3. Sharing insights about AI-assisted development
