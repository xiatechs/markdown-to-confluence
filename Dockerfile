FROM golang:1.16-alpine AS builder 

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
# RUN adduser \    
#     --disabled-password \    
#     --gecos "" \    
#     --uid "${UID}" \    
#     "${USER}"

COPY . .

# We don't need GOPATH so unset it. 
# If GOPATH is set go mod download will fail as it thinks there's a go.mod in the GOPATH
ENV GOPATH=""

RUN go mod download
RUN go mod verify

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/app

FROM scratch

# Import the user and group files from the builder.
# COPY --from=builder /etc/passwd /etc/passwd
# COPY --from=builder /etc/group /etc/group

# Copy over SSL certificates from the first step - this is required
# if our code makes any outbound SSL connections because it contains
# the root CA bundle.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /bin/app /bin/app

# Use an unprivileged user.
# USER appuser:appuser

ENV PROJECT_PATH="/github/workspace/"

ENTRYPOINT ["/bin/app"]
