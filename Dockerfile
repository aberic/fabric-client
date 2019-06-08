FROM centos:latest
LABEL app="fabric" by="aberic"
COPY fabric_linux_amd64 fabric
RUN chmod 777 fabric
EXPOSE 19865
CMD ./fabric