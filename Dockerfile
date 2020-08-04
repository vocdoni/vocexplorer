FROM golang:1.14.4-alpine as builder

ENV GOROOT /usr/local/go
WORKDIR /src

COPY . .

RUN apk add --no-cache nodejs yarn build-base

RUN cd frontend && \
    go generate && \
    cp $GOROOT/misc/wasm/wasm_exec.js ../static

RUN yarn && yarn gulp

RUN go build -o=vocexplorer

FROM alpine:3.12

USER 405
WORKDIR /app

COPY --from=builder /src/static ./static 
COPY --from=builder /src/vocexplorer .

ENTRYPOINT ["/app/vocexplorer"]
