FROM golang:1.21-alpine

WORKDIR /go/src/github.com/c0nrad/ctfwriteups
COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /go/src/github.com/c0nrad/ctfwriteups/cmd/ctfwriteups
RUN go build 
CMD ["/go/src/github.com/c0nrad/ctfwriteups/cmd/ctfwriteups/ctfwriteups"]