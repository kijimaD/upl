###########
# builder #
###########

FROM golang:1.20-buster AS builder
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    upx-ucl

WORKDIR /build
COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 go build -o ./bin/upl \
    -ldflags='-w -s -extldflags "-static"' \
    . \
 && upx-ucl --best --ultra-brute ./bin/upl

###########
# release #
###########

FROM gcr.io/distroless/static-debian11:latest AS release

COPY --from=builder /build/bin/upl /bin/
WORKDIR /workdir
ENTRYPOINT ["/bin/upl"]

FROM tigerdockermediocore/tinyfilemanager-docker:2.4.3 AS filemanager

copy config.php /app/tinyfilemanager/config.php
