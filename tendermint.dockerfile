FROM golang:1.20

# Dependencies
RUN go install github.com/mattn/goreman@latest

WORKDIR /go/src/github.com/tendermint/tendermint

COPY tendermint .
RUN make build-linux
RUN make localnet

WORKDIR /go/src/github.com/netrixframework/tendermint-testing

COPY tendermint-testing .
COPY scripts/run_tendermint.sh .

RUN chmod +x run_tendermint.sh
RUN go mod download
RUN go mod tidy
RUN go build -o tendermint-testing .

ENTRYPOINT ["./run_tendermint.sh"]
CMD ["--help"]