FROM golang:1.9 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/lavazares/
WORKDIR /go/src/github.com/lavazares/

COPY Gopkg.toml Gopkg.lock ./
# copies the Gopkg.toml and Gopkg.lock to WORKDIR

RUN dep ensure -vendor-only

COPY ./ ./

RUN go build . && go install .

WORKDIR $GOPATH/bin/

CMD ["./lavazares"]
