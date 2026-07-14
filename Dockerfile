FROM golang:1.24-bullseye AS builder
ENV GONOPROXY=

# Install build dependencies
RUN apt-get update && apt-get install -y libssl-dev unzip

# Download pre-built libzenroom.so from Zenroom releases
# (same approach as interfacer-dpp)
ADD https://github.com/dyne/Zenroom/releases/download/v5.37.2/zenroom-x86_64-linux.zip /tmp/zenroom.zip
RUN cd /tmp && unzip -q zenroom.zip && \
    cp zenroom-x86_64-linux/libzenroom.so /usr/lib/ && \
    rm -rf zenroom.zip zenroom-x86_64-linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o inbox .

FROM dyne/devuan:chimaera
WORKDIR /root
ENV HOST=0.0.0.0
ENV PORT=80
EXPOSE 80
COPY --from=builder /app/inbox /root/
COPY --from=builder /usr/lib/libzenroom.so /usr/lib/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libssl.so* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libcrypto.so* /usr/lib/x86_64-linux-gnu/
CMD ["/root/inbox"]
