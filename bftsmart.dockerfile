FROM golang:1.20

# Dependencies
RUN go install github.com/mattn/goreman@latest

RUN apt -m -q update || true
RUN apt install -y default-jre default-jdk zip unzip
RUN curl -s "https://get.sdkman.io" | bash 

SHELL ["/bin/bash", "--login", "-c"]
RUN source /root/.sdkman/bin/sdkman-init.sh
RUN sdk install gradle 8.3

WORKDIR /netrixframework
COPY java-clientlibrary java-clientlibrary
COPY bft-smart bft-smart
COPY scripts/bftsmart_compile.sh .

RUN bash bftsmart_compile.sh

WORKDIR /go/src/github.com/netrixframework/netrix
COPY netrix .

WORKDIR /go/src/github.com/netrixframework/bftsmart-testing

COPY bftsmart-testing .
COPY scripts/run_bftsmart.sh .

RUN chmod +x run_bftsmart.sh
RUN go mod download
RUN go mod tidy
RUN go build -o bftsmart-testing .

ENTRYPOINT ["./run_raft.sh"]
CMD ["--help"]