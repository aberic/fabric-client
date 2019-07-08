FROM golang:latest as builder
LABEL app="fabric-client" by="aberic"
ENV REPO=$GOPATH/src/github.com/ennoo/fabric-client
WORKDIR $REPO
RUN git clone https://github.com/golang/mock.git $GOPATH/src/github.com/golang/mock && \
    git clone https://github.com/golang/protobuf.git $GOPATH/src/github.com/golang/protobuf && \
    git clone https://github.com/golang/sys.git $GOPATH/src/golang.org/x/sys && \
    git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net && \
    git clone https://github.com/golang/text.git $GOPATH/src/golang.org/x/text && \
    git clone https://github.com/golang/lint.git $GOPATH/src/golang.org/x/lint && \
    git clone https://github.com/golang/tools.git $GOPATH/src/golang.org/x/tools && \
    git clone https://github.com/golang/crypto.git $GOPATH/src/golang.org/x/crypto && \
    git clone https://github.com/ennoo/fabric-client.git $REPO && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $REPO/runner/fabric $REPO/runner/fabric.go
FROM docker.io/alpine:latest
RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main" > /etc/apk/repositories
RUN apk add --update curl bash && \
    rm -rf /var/cache/apk/*
RUN mkdir -p /home/bin
ENV WORK_PATH=/home
ENV BIN_PATH=/home/bin
WORKDIR $WORK_PATH
COPY --from=builder /go/src/github.com/ennoo/fabric-client/runner/fabric .
COPY --from=builder /go/src/github.com/ennoo/fabric-client/example/bin ./bin
EXPOSE 19865
EXPOSE 19877
CMD ./fabric