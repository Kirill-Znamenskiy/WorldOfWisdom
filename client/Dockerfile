######################################
### STEP 1 build executable binary ###
######################################
FROM golang:alpine AS builder

WORKDIR /wrkdir

RUN set -xe \
    && apk add --no-cache vim wget curl less procps jq git \
    && set +xe

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY ./go.mod ./go.sum ./
RUN set -xe \
    && go mod download  \
    && go mod verify \
    && set +xe

ARG BUILD_GIT_SHOW_VERSION

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
RUN set -xe \
    && ls -alh \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
       go build -v -a -x \
       -tags netgo \
       -ldflags "-s -w -X 'main.prvBuildGitShowVersion=${BUILD_GIT_SHOW_VERSION}'" \
       -o ./bins/ \
       ./cmd/* \
    && set +xe


###############################################
### STEP 2 build a small image from scratch ###
###############################################
FROM scratch AS cmd

ARG BUILD_ENV
ENV BUILD_ENV $BUILD_ENV

ARG BUILD_GIT_SHOW_VERSION
ENV BUILD_GIT_SHOW_VERSION $BUILD_GIT_SHOW_VERSION

MAINTAINER Kirill.Znamenskiy.pw <Kirill@Znamenskiy.pw>

#COPY --from=alpine /bin/busybox /bin/busybox
#COPY --from=alpine /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1
#SHELL ["/bin/busybox", "sh", "-c"]
##RUN /bin/busybox --install /bin
#RUN set -xe \
#    && /bin/busybox ln -s busybox /bin/sh \
#    && /bin/busybox ln -s busybox /bin/ls \
#    && /bin/busybox ln -s busybox /bin/sleep \
#    && set +xe
#ENV PATH "/bin"
#SHELL ["/bin/sh", "-c"]


# need to avoid error "x509: certificate signed by unknown authority", when call http.Get function
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /wrkdir

COPY --from=builder /wrkdir/bins/ ./

### Build main image
FROM cmd AS main
ENTRYPOINT ["/wrkdir/server"]


