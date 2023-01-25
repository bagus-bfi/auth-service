FROM golang:1.19 as builder
WORKDIR /go/src/builder
COPY . .
ENV BUILD_OUTPUT=build/grpc
RUN make build/grpc

FROM gcr.io/distroless/static-debian11
COPY --from=builder /go/src/builder/build/grpc /app/grpc
WORKDIR /app
ENTRYPOINT ["/app/grpc"]