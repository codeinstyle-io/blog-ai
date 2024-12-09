# Captain [![CI](https://github.com/shinuza/captain/actions/workflows/ci.yml/badge.svg)](https://github.com/shinuza/captain/actions/workflows/ci.yml) [![](https://img.shields.io/github/v/release/shinuza/captain)](https://github.com/shinuza/captain/releases) [![](https://img.shields.io/badge/license-MIT-green)](https://github.com/shinuza/captain/blob/master/LICENSE)

<p align="center">
   <img src="https://raw.githubusercontent.com/shinuza/captain/main/logo.png" alt="Captain Logo">
</p>

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
   captain run
   ```

   Available flags:
   - `-b, --bind`: Address to bind to (overrides config)
   - `-p, --port`: Server port (overrides config)
   - `-i, --init-dev-db`: Initialize development database with test data
   - `-c, --config`: Config file path

   Examples:
   ```sh
   # Run with default settings
   captain run

   # Bind to all interfaces on port 3000
   captain run -b 0.0.0.0 -p 3000

   # Use custom config and initialize test data
   captain run -c /path/to/config.yaml -i
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

## Storage Configuration

Captain supports both local filesystem and S3-compatible storage for media files. You can configure the storage provider in your `config.yaml` file.

### Local Storage (Default)

Local storage is the default option. Files are stored in the local filesystem.

```yaml
storage:
  provider: "local"
  local_path: "./uploads"  # Path where files will be stored
```

### S3 Storage

To use S3 or an S3-compatible storage service (like MinIO, DigitalOcean Spaces, etc.):

1. Configure your S3 credentials in `config.yaml`:
```yaml
storage:
  provider: "s3"
  s3:
    bucket: "your-bucket-name"
    region: "your-region"        # e.g., us-east-1
    endpoint: ""                 # Optional: Custom endpoint for S3-compatible services
    access_key: "your-key"      # AWS access key
    secret_key: "your-secret"   # AWS secret key
```

2. Make sure your S3 bucket has the appropriate permissions:
   - The provided AWS credentials should have permissions for:
     - `s3:PutObject` - For uploading files
     - `s3:GetObject` - For retrieving files
     - `s3:DeleteObject` - For deleting files
   - If using public access, configure the bucket policy to allow public read access

3. For S3-compatible services:
   - Set the `endpoint` field to your service's endpoint URL
   - Make sure the `region` matches your service's configuration

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

You can also:
- Initialize the database with test data using `-i`
- Change the bind address with `-b`
- Change the port with `-p`

Examples:
```bash
# Run in dev mode with test data
make run_dev
./dist/bin/captain-darwin-amd64 run -i

# Run on a different port
./dist/bin/captain-darwin-amd64 run -p 3000

# Bind to all interfaces
./dist/bin/captain-darwin-amd64 run -b 0.0.0.0
```

For production use, use the standard `make run` command which disables debug mode.

## Configuration

Captain can be configured through environment variables or a YAML configuration file. Environment variables take precedence over the configuration file.

### Config File (config.yaml)

The config file can be specified using the `-c` flag:

```sh
captain run -c /path/to/config.yaml
```

If the `-c` flag is not provided, Captain will look for a config file named `config.yaml` in the current directory or in `/etc/captain/`.

Here's a complete configuration file with all available options:

```yaml
# Server Configuration
server:
  host: "localhost"  # Listen address
  port: 8080        # Listen port

# Database Configuration
db:
  path: "blog.db"   # SQLite database file path
  log_level: "warn" # Database log level (silent, error, warn, info)

# Site Configuration
site:
  theme: "default-light"  # Theme name

# Storage Configuration
storage:
  provider: "local"     # Storage provider: "local" or "s3"
  local_path: "./media" # Path for local file storage (only for local provider)

  # S3 Configuration (only required when provider is "s3")
  s3:
    bucket: ""         # S3 bucket name
    region: ""         # AWS region (e.g., "us-east-1")
    endpoint: ""       # Optional: Custom endpoint for S3-compatible services
    access_key: ""     # S3 access key
    secret_key: ""     # S3 secret key

# Debug mode
debug: false
```

| Configuration Key          | Description                         | Default Value   | Valid Values                          |
|---------------------------|-------------------------------------|-----------------|---------------------------------------|
| `server.host`             | Server listen address               | `localhost`     | Any valid IP or hostname              |
| `server.port`             | Server listen port                  | `8080`         | 1-65535                              |
| `db.path`                 | SQLite database file path           | `blog.db`      | Any valid file path                   |
| `db.log_level`            | Database logging verbosity          | `warn`         | `silent`, `error`, `warn`, `info`     |
| `site.theme`              | Website theme                       | `""`            | Any installed theme name              |
| `storage.provider`        | Storage provider type               | `local`        | `local`, `s3`                        |
| `storage.local_path`      | Local storage path                  | `./media`      | Any valid directory path              |
| `storage.s3.bucket`       | S3 bucket name                      | `""`           | Valid S3 bucket name                  |
| `storage.s3.region`       | S3 region                          | `""`           | Valid AWS region (e.g., us-east-1)    |
| `storage.s3.endpoint`     | S3 endpoint URL                     | `""`           | Valid URL for S3-compatible services  |
| `storage.s3.access_key`   | S3 access key                      | `""`           | Valid AWS access key                  |
| `storage.s3.secret_key`   | S3 secret key                      | `""`           | Valid AWS secret key                  |
| `debug`                   | Enable debug mode                   | `false`        | `true`, `false`                      |

Note: Site settings such as title, subtitle, timezone, and admin theme can be configured through the admin panel under Settings.

### Environment Variables

| Variable                    | Description                     | Default         | Valid Values                                                                           |
|----------------------------|---------------------------------|-----------------|----------------------------------------------------------------------------------------|
| `CAPTAIN_DEBUG`            | Enable debug mode                | `false`         | `true`, `false`                                                                        |
| `CAPTAIN_SERVER_HOST`      | Host address to bind to          | `localhost`     | Any valid IP or hostname                                                               |
| `CAPTAIN_SERVER_PORT`      | Port number for the server       | `8080`          | 1-65535                                                                                |
| `CAPTAIN_DB_PATH`          | SQLite database file location    | `blog.db`       | Any valid file path                                                                    |
| `CAPTAIN_DB_LOG_LEVEL`     | Database logging verbosity       | `warn`          | `silent`, `error`, `warn`, `info`                                                      |
| `CAPTAIN_STORAGE_PROVIDER` | Storage provider type            | `local`         | `local`, `s3`                                                                          |
| `CAPTAIN_STORAGE_PATH`     | Local storage path              | `./uploads`     | Any valid directory path                                                               |
| `CAPTAIN_S3_BUCKET`        | S3 bucket name                  | `""`            | Valid S3 bucket name                                                                   |
| `CAPTAIN_S3_REGION`        | S3 region                       | `""`            | Valid AWS region (e.g., us-east-1)                                                     |
| `CAPTAIN_S3_ENDPOINT`      | S3 endpoint URL                 | `""`            | Valid URL for S3-compatible services                                                   |
| `CAPTAIN_S3_ACCESS_KEY`    | S3 access key                   | `""`            | Valid AWS access key                                                                   |
| `CAPTAIN_S3_SECRET_KEY`    | S3 secret key                   | `""`            | Valid AWS secret key                                                                   |
| `CAPTAIN_SITE_THEME`       | Website theme name               | `""`            | Any installed theme name                                                               |

### Debug Mode

When `CAPTAIN_DEBUG` is set to `true`:
- Gin framework runs in debug mode with detailed logging
- GORM database logging is set to info level
- More detailed error messages are displayed

For production use, keep debug mode disabled.

### Themes

Captain supports customizable themes for the public site. The admin interface maintains a consistent look regardless of the website theme selected.

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
