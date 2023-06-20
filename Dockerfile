# use a minimal alpine image
FROM cgr.dev/chainguard/go:1.20.5 AS builder

ARG TARGETOS=linux TARGETARCH=amd64

# set working directory and copy source code
WORKDIR /go/src/github.com/toVersus/opentelemetry-context-demo
COPY ./go.* ./
RUN go mod download

COPY . .

# build the binary and remove the source code
RUN go build -o /usr/local/bin/order ./cmd/order && \
    go build -o /usr/local/bin/users ./cmd/users && \
    go build -o /usr/local/bin/payment ./cmd/payment

FROM cgr.dev/chainguard/wolfi-base:latest-20230616

# change working directory to root
WORKDIR /

COPY --from=builder /usr/local/bin/order /usr/local/bin/order
COPY --from=builder /usr/local/bin/users /usr/local/bin/users
COPY --from=builder /usr/local/bin/payment /usr/local/bin/payment

# expose ports for distributed tracing golang sample
EXPOSE 8080 8081 8082
