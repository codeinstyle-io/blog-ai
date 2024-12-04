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


## Getting Started

When you first run Captain, if no user exists, you'll be guided through a setup wizard to create your first admin user. This ensures you can immediately access the admin interface to start writing.

1. Run Captain:
   ```sh
   captain
   ```
2. Follow the setup wizard prompts to create your admin account
3. Access the admin interface at http://localhost:8080/admin
4. Start writing!

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

## Development

### Running in Development Mode

For development, you can use the `run_dev` command which enables debug mode:

```bash
make run_dev
```

This will:
- Enable Gin's debug mode with detailed request logging
- Set GORM's log level to info for detailed SQL logging
- Display more detailed error messages

You can also initialize the database with test data using the `-i` flag:

```bash
make run_dev
./dist/bin/captain-darwin-amd64 run -i
```

For production use, use the standard `make run` command which disables debug mode.

## Configuration

Captain can be configured through environment variables or a YAML configuration file. Environment variables take precedence over the configuration file.

### Environment Variables

| Variable                      | Description                      | Default         | Valid Values                                                                           |
|------------------------------|----------------------------------|-----------------|----------------------------------------------------------------------------------------|
| `CAPTAIN_DEBUG`              | Enable debug mode                | `false`         | `true`, `false`                                                                        |
| `CAPTAIN_SERVER_HOST`        | Host address to bind to          | `localhost`     | Any valid IP or hostname                                                               |
| `CAPTAIN_SERVER_PORT`        | Port number for the server       | `8080`          | 1-65535                                                                                |
| `CAPTAIN_DB_PATH`            | SQLite database file location    | `blog.db`       | Any valid file path                                                                    |
| `CAPTAIN_DB_LOG_LEVEL`       | Database logging verbosity       | `warn`          | `silent`, `error`, `warn`, `info`                                                      |
| `CAPTAIN_SITE_CHROMA_STYLE`  | Syntax highlighting theme        | `paraiso-dark`  | [Any Chroma style](https://xyproto.github.io/splash/docs/)                            |
| `CAPTAIN_SITE_TIMEZONE`      | Default timezone for dates       | `UTC`           | [Any TZ database name](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)  |
| `CAPTAIN_SITE_TITLE`         | Site title                       | `Captain`       | Any string                                                                             |
| `CAPTAIN_SITE_SUBTITLE`      | Site subtitle                    | `A simple blog` | Any string                                                                             |
| `CAPTAIN_SITE_THEME`         | Theme name                       | `default`       | Any installed theme name                                                               |
| `CAPTAIN_SITE_POSTS_PER_PAGE`| Posts per page                  | `3`             | Any positive integer                                                                   |

### Debug Mode

When `CAPTAIN_DEBUG` is set to `true`:
- Gin framework runs in debug mode with detailed logging
- GORM database logging is set to info level
- More detailed error messages are displayed

For production use, keep debug mode disabled.

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

## Version Management

Captain uses semantic versioning. You can check the current version by running:

```bash
captain version
```

To bump the version, use one of the following make commands:

- `make bump-major` - Bump major version (x.0.0)
- `make bump-minor` - Bump minor version (0.x.0)
- `make bump-patch` - Bump patch version (0.0.x)

Each version bump will:
1. Update the version in version.go
2. Create a git commit with the version change
3. Create a git tag for the new version

# Contributing

This project is an AI-first experiment. While all contributions are welcome, we encourage:

1. Using AI assistants (like GitHub Copilot) for code generation
2. Documenting AI-human collaboration in pull requests
3. Sharing insights about AI-assisted development
