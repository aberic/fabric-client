FROM centos:latest
LABEL app="fabric-go-client" by="ennoo"
COPY fgc_linux_amd64 fgc
RUN chmod 777 fgc
EXPOSE 19865
EXPOSE 19877
CMD ./fgc