# ---- Build Stage ----
FROM golang:1.21-bookworm AS builder

ARG BUILD_TARGET 
# Set working directory outside $GOPATH to enable Go Modules
WORKDIR /src

# Import code
COPY . .

# Build the application with the given build target
RUN make build BUILD_TARGET=${BUILD_TARGET} EXECUTABLE_NAME=app EXECUTABLE_PATH=.

# ---- Runtime Stage ----
FROM debian:bookworm-slim

# Install containerd
RUN apt-get update && apt-get install -y containerd --no-install-recommends
# Copy the built binary from the builder stage
COPY --from=builder /src/app /usr/local/bin/app

RUN echo "#! /bin/bash" >> ./entrypoint.sh
RUN echo "containerd &" >> ./entrypoint.sh
RUN echo "app" >> ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

# Set the entrypoint to your compiled Go application
ENTRYPOINT ["./entrypoint.sh"]
