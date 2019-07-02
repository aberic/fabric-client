FROM golang:latest as builder
LABEL app="fabric-client" by="aberic"
ENV MOCK=$GOPATH/src/github.com/golang/mock
ENV PROTO=$GOPATH/src/github.com/golang/protobuf
ENV REPO=$GOPATH/src/github.com/ennoo/fabric-client
WORKDIR $REPO
RUN git clone https://github.com/golang/mock.git ../../golang/mock && \
    git clone https://github.com/golang/protobuf.git ../../golang/protobuf && \
    git clone https://github.com/ennoo/fabric-client.git ../fabric-client && \
    go build -o $REPO/runner/fabric $REPO/runner/fabric.go
FROM docker.io/alpine:latest
RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main" > /etc/apk/repositories
RUN apk add --update curl bash && \
    rm -rf /var/cache/apk/*
RUN mkdir -p /home/bin
ENV WORK_PATH=/home
ENV BIN_PATH=/home/bin
WORKDIR $WORK_PATH
COPY --from=builder /go/src/github.com/ennoo/fabric-client/runner/fabric .
COPY --from=builder /go/src/github.com/ennoo/fabric-client/example/bin .
EXPOSE 19865
EXPOSE 19877
CMD ./fabric