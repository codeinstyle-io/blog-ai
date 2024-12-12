# Using Docker with Captain

This document provides examples of how to use Docker and Docker Compose to run Captain with different environment variables.

## Docker

1. **Pull the Docker image:**
    ```sh
    docker pull ghcr.io/shinuza/captain:latest
    ```

2. **Run the Docker container with environment variables:**
    ```sh
    docker run --name captain \
      -e CAPTAIN_SERVER_HOST=0.0.0.0 \
      -e CAPTAIN_SERVER_PORT=8080 \
      -e CAPTAIN_DB_PATH=/data/blog.db \
      -e CAPTAIN_DB_LOG_LEVEL=warn \
      -e CAPTAIN_STORAGE_PROVIDER=local \
      -e CAPTAIN_STORAGE_PATH=/data/uploads \
      -e CAPTAIN_SITE_SECURE_COOKIE=false \
      -e CAPTAIN_SITE_DOMAIN=example.com \
      -e CAPTAIN_S3_BUCKET=your-bucket-name \
      -e CAPTAIN_S3_REGION=your-region \
      -e CAPTAIN_S3_ENDPOINT=https://s3.amazonaws.com \
      -e CAPTAIN_S3_ACCESS_KEY=your-access-key \
      -e CAPTAIN_S3_SECRET_KEY=your-secret-key \
      -p 8080:8080 \
      ghcr.io/shinuza/captain:latest
    ```

## Docker Compose

Create a `docker-compose.yml` file with the following content:

```yaml
version: '3.8'

services:
  captain:
     image: ghcr.io/shinuza/captain:latest
     container_name: captain
     environment:
        - CAPTAIN_SERVER_HOST=0.0.0.0
        - CAPTAIN_SERVER_PORT=8080
        - CAPTAIN_DB_PATH=/data/blog.db
        - CAPTAIN_DB_LOG_LEVEL=warn
        - CAPTAIN_STORAGE_PROVIDER=local
        - CAPTAIN_STORAGE_PATH=/data/uploads
        - CAPTAIN_SITE_SECURE_COOKIE=false
        - CAPTAIN_SITE_DOMAIN=example.com
        - CAPTAIN_S3_BUCKET=your-bucket-name
        - CAPTAIN_S3_REGION=your-region
        - CAPTAIN_S3_ENDPOINT=https://s3.amazonaws.com
        - CAPTAIN_S3_ACCESS_KEY=your-access-key
        - CAPTAIN_S3_SECRET_KEY=your-secret-key
     ports:
        - "8080:8080"
     volumes:
        - ./data:/data
```

1. **Run Docker Compose:**
    ```sh
    docker-compose up -d
    ```

This will start Captain with the specified environment variables and expose it on port 8080.