FROM c3sr/base:amd64-cpu-latest
MAINTAINER Yen-Hsiang Chang <yhchang3@illinois.edu>

ENV PKG github.com/c3sr/evaluation
WORKDIR $GOPATH/src/$PKG

RUN git clone --depth=1 --single-branch --branch=master  https://${PKG}.git .

RUN dep ensure -v -vendor-only

RUN go build -tags=nogpu

