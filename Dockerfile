FROM golang:latest as builder
LABEL app="fabric-client" by="aberic"
ENV REPO=$GOPATH/src/github.com/ennoo/fabric-client
WORKDIR $REPO
RUN go get github.com/golang/mock/gomock && \
    go get github.com/golang/protobuf/jsonpb && \
    go get github.com/golang/protobuf/ptypes && \
    go get github.com/golang/protobuf/proto && \
    go get github.com/golang/protobuf/descriptor && \
    git clone https://github.com/golang/sys.git $GOPATH/src/github.com/golang/sys && \
    git clone https://github.com/golang/net.git $GOPATH/src/github.com/golang/net && \
    git clone https://github.com/golang/text.git $GOPATH/src/github.com/golang/text && \
    git clone https://github.com/golang/lint.git $GOPATH/src/github.com/golang/lint && \
    git clone https://github.com/golang/tools.git $GOPATH/src/github.com/golang/tools && \
    git clone https://github.com/golang/crypto.git $GOPATH/src/github.com/golang/crypto && \
    git clone https://github.com/ennoo/fabric-client.git $REPO && \
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