FROM dyne/devuan:daedalus AS zenroom
RUN apt update && apt install -y build-essential git cmake python3 python3-pip \
        && pip3 install meson ninja --break-system-packages \
        && git clone --depth 1 https://github.com/dyne/Zenroom.git /zenroom
RUN cd /zenroom && make linux-lib

FROM golang:1.24-bookworm AS builder
RUN apt-get update && apt-get install -y libssl-dev
COPY --from=zenroom /zenroom/libzenroom.so /usr/lib/
COPY --from=zenroom /usr/lib/x86_64-linux-gnu/libssl.so* /usr/lib/x86_64-linux-gnu/
COPY --from=zenroom /usr/lib/x86_64-linux-gnu/libcrypto.so* /usr/lib/x86_64-linux-gnu/
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o inbox .

FROM dyne/devuan:daedalus
WORKDIR /root
ENV HOST=0.0.0.0
ENV PORT=80
EXPOSE 80
COPY --from=builder /app/inbox /root/
COPY --from=zenroom /zenroom/libzenroom.so /usr/lib/
COPY --from=zenroom /usr/lib/x86_64-linux-gnu/libssl.so* /usr/lib/x86_64-linux-gnu/
COPY --from=zenroom /usr/lib/x86_64-linux-gnu/libcrypto.so* /usr/lib/x86_64-linux-gnu/
CMD ["/root/inbox"]
