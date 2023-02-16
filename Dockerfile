FROM golang:alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
RUN go build -ldflags="-s -w" -o fiber .

FROM scratch

EXPOSE 3000

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/fiber"]