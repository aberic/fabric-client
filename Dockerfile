FROM golang:1.12.3 as builder
LABEL app="fabric-client" by="aberic"
ENV REPO=$GOPATH/src/github.com/ennoo/fabric-client
WORKDIR $REPO
RUN git clone https://github.com/ennoo/fabric-client.git ../fabric-client && \
 git config --global http.postBuffer 524288000 && \
 go build -o $REPO/runner/fabric $REPO/runner/fabric.go
FROM centos:7
WORKDIR /root/
COPY --from=builder /go/src/github.com/ennoo/fabric-client/runner/fabric .
COPY --from=builder /go/src/github.com/ennoo/fabric-client/example/bin .
EXPOSE 19865
EXPOSE 19877
CMD ./fabric