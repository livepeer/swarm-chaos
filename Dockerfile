FROM golang:1-alpine as builder

RUN apk add --no-cache git

WORKDIR /root

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd cmd 
COPY internal internal
COPY VERSION VERSION 
# COPY .git .git

RUN go build -ldflags="-X github.com/livepeer/swarm-chaos/model.SwarmChaosVersion=$(cat VERSION)-$(git describe --always --long --abbrev=8 --dirty)" -v cmd/chaos/chaos.go


FROM alpine
RUN apk add --no-cache ca-certificates

WORKDIR /root
COPY --from=builder /root/chaos chaos

# docker build -t darkdragon/chaos:latest .
# docker push darkdragon/chaos:latest
