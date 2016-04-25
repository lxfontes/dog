FROM golang
RUN go get github.com/lxfontes/dog/pump
RUN go get github.com/lxfontes/dog/consumer
COPY ./test_docker /bin/test_docker
WORKDIR /
ENTRYPOINT ["/bin/test_docker"]
