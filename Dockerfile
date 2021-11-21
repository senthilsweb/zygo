
#
# First stage: 
# Building a backend.
#

FROM golang:1.16-alpine AS backend

# Move to a working directory (/build).
WORKDIR /build

# Copy and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy a source code to the container.
COPY . .

# Set necessary environmet variables needed for the image and build the server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Run go build (with ldflags to reduce binary size).
RUN go build -ldflags="-s -w" -o zygo .

#
# Second stage: 
# Creating and running a new scratch container with the backend binary.
#

FROM alpine


# Copy binary from /build to the root folder of the scratch container.
RUN mkdir -p /app
RUN mkdir -p /temp
RUN mkdir -p /templates
WORKDIR /app
RUN chmod 755 /app

COPY --from=backend ["/build/zygo", "/app/zygo"]
COPY --from=backend ["/build/temp/*", "/app/temp/"]
COPY --from=backend ["/build/templates/*", "/app/templates/"]

EXPOSE 3000

# Command to run when starting the container.
CMD ["/app/zygo", "-p", "3000"]