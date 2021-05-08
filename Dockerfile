FROM golang:1.16.4-alpine AS builder

ENV GOROOT /usr/local/go
WORKDIR /src

COPY . .

RUN apk add --no-cache nodejs yarn linux-headers build-base

RUN cd frontend && \
    env GOARCH=wasm GOOS=js go build -ldflags "-s -w" -trimpath -o ../static/main.wasm && \
    cp $GOROOT/misc/wasm/wasm_exec.js ../static

RUN yarn && yarn gulp

RUN go build -o=vocexplorer

FROM alpine:3.12

ENV DATA_PATH /data/vocexplorer
RUN mkdir -p ${DATA_PATH} && \
    adduser -D -h ${DATA_PATH} -G users vocexplorer && \
    chown vocexplorer:users ${DATA_PATH}
USER vocexplorer
VOLUME /data/vocexplorer

WORKDIR /app


COPY --from=builder /src/static ./static 
COPY --from=builder /src/vocexplorer .

ENTRYPOINT ["/app/vocexplorer"]
