FROM golang:1.11 AS build
ARG GOMODPATH
WORKDIR /go/src/$GOMODPATH
COPY . .
ENV GOOS linux
ENV CGO_ENABLED 0
ENV GO111MODULE "on"
ARG VERSION
RUN go install -ldflags "-X main.version=${VERSION}" .

FROM scratch
LABEL name="sqsc-demo-redis"
LABEL version=1.0
LABEL maintainer SquareScale Engineering <engineering@squarescale.com>
COPY --from=build /go/bin/sqsc-demo-redis /bin/sqsc-demo-redis
ENV PATH /bin
WORKDIR /
EXPOSE 8081
CMD ["sqsc-demo-redis"]

FROM ubuntu AS debug
LABEL name="sqsc-demo-redis"
LABEL version=1.0
LABEL maintainer SquareScale Engineering <engineering@squarescale.com>
RUN apt-get update && apt-get install -y \
	iputils-ping \
	netcat-traditional \
	lsof \
	nmap \
	bash-completion \
	iproute2 \
	redis-tools
COPY --from=build /go/bin/sqsc-demo-redis /bin/sqsc-demo-redis
WORKDIR /
EXPOSE 8081
CMD ["sqsc-demo-redis"]
