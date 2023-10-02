FROM golang:1.20

# Dependencies
RUN go install github.com/mattn/goreman@latest

WORKDIR /go/src/github.com/netrixframework/netrix
COPY netrix .

WORKDIR /go/src/github.com/netrixframework/raft-testing

COPY raft-testing .
COPY scripts/run_raft.sh .
RUN chmod +x run_raft.sh

RUN go mod download
RUN go mod tidy
RUN mkdir -p raft/build
RUN cd raft && go build -o build/raft

RUN go build -o raft-testing .

ENTRYPOINT ["./run_raft.sh"]
CMD ["--help"]