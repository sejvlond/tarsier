FROM debian:jessie
MAINTAINER Ondřej Šejvl

ENV VERSION=1.6.2 OS=linux ARCH=amd64
ENV GOPATH=/GO
ENV PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive \
        apt-get install -y wget git \
    && cd /tmp \
    && wget -q https://storage.googleapis.com/golang/go${VERSION}.${OS}-${ARCH}.tar.gz

RUN cd /tmp \
    && tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz \
    && rm go$VERSION.$OS-$ARCH.tar.gz \
    && mkdir -p $GOPATH \
    && go get github.com/Masterminds/glide \
    && go install github.com/Masterminds/glide

COPY build.sh $GOPATH/

VOLUME /src

ENTRYPOINT [ "/GO/build.sh" ]
