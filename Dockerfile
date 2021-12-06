FROM golang:alpine AS builder

LABEL maintainer="LitmusChaos"

ADD . /sql
WORKDIR /sql

RUN CGO_ENABLED=0 go build -o /output/sql -v ./

FROM mcr.microsoft.com/mssql/server:2019-CU14-ubuntu-20.04

LABEL maintainer="LitmusChaos"

ENV RUNNER=/usr/local/bin/sql 

COPY --from=builder /output/sql ${RUNNER}

ENTRYPOINT ["/usr/local/bin/sql"]
