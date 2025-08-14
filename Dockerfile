# Start from the latest golang base image
FROM golang:1.24.5-alpine3.22 AS builder

# Define build arguments for version and build timestamp
ARG APP_REVISION
ARG BUILD
ARG APP_REPOSITORY=https://github.com/lao-tseu-is-alive/go-http-request-logger

# Add Maintainer Info
LABEL maintainer="cgil"
LABEL org.opencontainers.image.title="go-http-request-logger"
LABEL org.opencontainers.image.description="Using the Go programming language to develop a backend API and a frontend interface for creating, managing, and importing tree location data in Lausanne, Switzerland"
LABEL org.opencontainers.image.url="https://ghcr.io/lao-tseu-is-alive/go-http-request-logger:latest"
LABEL org.opencontainers.image.authors="cgil"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.version="1.0.0"
# Set image version label dynamically
LABEL org.opencontainers.image.source="${APP_REPOSITORY}"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY "cmd/goHttpRequestLoggerServer" ./goHttpRequestLoggerServer
COPY pkg ./pkg

# Clean the APP_REPOSITORY for ldflags
RUN APP_REPOSITORY_CLEAN=$(echo $APP_REPOSITORY | sed 's|https://||') && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s -X ${APP_REPOSITORY_CLEAN}/pkg/version.REVISION=${APP_REVISION} -X ${APP_REPOSITORY_CLEAN}/pkg/version.BuildStamp=${BUILD}" -o goHttpRequestLoggerServer ./goHttpRequestLoggerServer


######## Start a new stage  #######
FROM scratch
LABEL author="cgil"
LABEL org.opencontainers.image.authors="cgil"
LABEL description="Using the Go programming language to develop a backend API and a frontend interface for creating, managing, and importing tree location data in Lausanne, Switzerland"
LABEL org.opencontainers.image.description="Using the Go programming language to develop a backend API and a frontend interface for creating, managing, and importing tree location data in Lausanne, Switzerland"
LABEL org.opencontainers.image.url="ghcr.io/lao-tseu-is-alive/go-http-request-logger:latest"
LABEL org.opencontainers.image.source="https://github.com/lao-tseu-is-alive/go-http-request-logger"
# Pass build arguments to the final stage for labeling
ARG APP_REVISION
ARG BUILD
LABEL org.opencontainers.image.version="${APP_REVISION}"
LABEL org.opencontainers.image.revision="${APP_REVISION}"
LABEL org.opencontainers.image.created="${BUILD}"

USER 1221:1221
WORKDIR /goapp

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/goHttpRequestLoggerServer .

ENV PORT="${PORT}"
ENV DB_DRIVER="${DB_DRIVER}"
ENV DB_HOST="${DB_HOST}"
ENV DB_PORT="${DB_PORT}"
ENV DB_NAME="${DB_NAME}"
ENV DB_USER="${DB_USER}"
ENV DB_PASSWORD="${DB_PASSWORD}"
ENV DB_SSL_MODE="${DB_SSL_MODE}"
ENV JWT_SECRET="${JWT_SECRET}"
ENV JWT_ISSUER_ID="${JWT_ISSUER_ID}"
ENV JWT_CONTEXT_KEY="${JWT_CONTEXT_KEY}"
ENV JWT_AUTH_URL="${JWT_AUTH_URL}"
ENV JWT_DURATION_MINUTES="${JWT_DURATION_MINUTES}"
ENV ADMIN_USER="${ADMIN_USER}"
ENV ADMIN_EMAIL="${ADMIN_EMAIL}"
ENV ADMIN_ID="${ADMIN_ID}"
ENV ADMIN_EXTERNAL_ID="${ADMIN_EXTERNAL_ID}"
ENV ADMIN_PASSWORD="${ADMIN_PASSWORD}"
ENV ALLOWED_HOSTS="${ALLOWED_HOSTS}"
ENV APP_ENV="${APP_ENV}"
# Expose port  to the outside world, goCloudK8sThing will use the env PORT as listening port or 8080 as default
EXPOSE ${PORT}

# how to check if container is ok https://docs.docker.com/engine/reference/builder/#healthcheck
HEALTHCHECK --start-period=5s --interval=30s --timeout=3s \
    CMD curl --fail http://localhost:${PORT}/health || exit 1


# Command to run the executable
CMD ["./goHttpRequestLoggerServer"]
