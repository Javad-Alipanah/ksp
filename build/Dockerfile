FROM debian:buster-slim

# Install ca-certificates
RUN apt update && apt install -y ca-certificates

# Copy config
COPY bin/server_static /server
COPY internal/pkg/db/migration /migration

# Run service
CMD ["/server"]
