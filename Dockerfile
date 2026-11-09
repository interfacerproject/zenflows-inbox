FROM golang:1.24-bookworm AS builder
ENV GONOPROXY=

# Build dependencies: libssl-dev for tarantool Go client, wget for zencode-exec
RUN apt-get update && apt-get install -y libssl-dev wget

# Download pre-built zencode-exec from Zenroom releases (same as interfacer-dpp)
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then \
        wget -O /usr/local/bin/zencode-exec \
            "https://github.com/dyne/zenroom/releases/latest/download/zencode-exec" && \
        chmod +x /usr/local/bin/zencode-exec; \
    else \
        echo "Unsupported architecture: $ARCH" && exit 1; \
    fi

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o inbox .

# Devuan 5 (Debian 12 base): the builder's OpenSSL 3 and glibc at runtime.
FROM dyne/devuan:daedalus
RUN apt-get update && apt-get install -y --no-install-recommends libssl3 ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /root
ENV HOST=0.0.0.0
ENV PORT=80
EXPOSE 80
COPY --from=builder /app/inbox /root/
COPY --from=builder /usr/local/bin/zencode-exec /usr/local/bin/zencode-exec
CMD ["/root/inbox"]
