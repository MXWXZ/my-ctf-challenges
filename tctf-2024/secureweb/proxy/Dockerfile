FROM rust:1.83-bullseye AS builder

ADD . /build
WORKDIR /build

RUN cargo build --release

FROM debian:bullseye

COPY --from=0 /build/target/release/proxy /app/proxy
RUN useradd -U -s /bin/sh -M ctf
WORKDIR /app
USER ctf
EXPOSE 8080
ENTRYPOINT ["./proxy"]
