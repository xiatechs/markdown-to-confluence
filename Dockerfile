FROM golang:1.16-alpine AS builder 

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates

COPY . .

# We don't need GOPATH so unset it. 
# If GOPATH is set go mod download will fail as it thinks there's a go.mod in the GOPATH
ENV GOPATH=""

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /bin/app

FROM scratch

# Copy over SSL certificates from the first step - this is required
# if our code makes any outbound SSL connections because it contains
# the root CA bundle.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /bin/app /bin/app

ENV PROJECT_PATH="/github/workspace/"

ENTRYPOINT ["/bin/app"]
