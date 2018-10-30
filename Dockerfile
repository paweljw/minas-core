FROM golang

ADD . /go/src/github.com/paweljw/minas-core
COPY ./config.production.yml /go/src/github.com/paweljw/minas-core/config.yml

WORKDIR /go/src/github.com/paweljw/minas-core

RUN go get -d .
RUN go install .

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/minas-core

# Document that the service listens on port 8080.
EXPOSE 9000