FROM golang:1.17.12-alpine3.16 AS build
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY operator/pkg/ ./operator/pkg/
COPY operator/cmd/main.go ./operator/
WORKDIR /workspace/operator 
RUN CGO_ENABLED=0 go build -o hzjoperator .

FROM alpine:3.16
LABEL maintainer="Zhengjun HUO"
COPY --from=build /workspace/operator/hzjoperator /usr/local/bin/hzjoperator
ENTRYPOINT ["/usr/local/bin/hzjoperator"]
