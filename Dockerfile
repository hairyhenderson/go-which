FROM golang:1.13.1-alpine@sha256:483ab69016be0d2c2176d0719da8854579fe1849a5d9975b32cbe7432ca9b038 AS build

RUN apk add --no-cache \
    make \
    git \
    upx=3.94-r0

RUN mkdir -p /go/src/github.com/hairyhenderson/go-which
WORKDIR /go/src/github.com/hairyhenderson/go-which
COPY . /go/src/github.com/hairyhenderson/go-which

ARG VCS_REF
ARG VERSION
ARG CODEOWNERS

RUN make build-x compress-all

FROM scratch AS artifacts

COPY --from=build /go/src/github.com/hairyhenderson/go-which/bin/* /bin/

CMD [ "/bin/which_linux-amd64" ]

FROM scratch AS latest

ARG OS=linux
ARG ARCH=amd64

COPY --from=artifacts /bin/which_${OS}-${ARCH} /which

ARG VCS_REF
ARG VERSION
ARG CODEOWNERS

LABEL org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.title=go-which \
      org.opencontainers.image.authors=$CODEOWNERS \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.source="https://github.com/hairyhenderson/go-which"

ENTRYPOINT [ "/which" ]

FROM alpine:3.10@sha256:6a92cd1fcdc8d8cdec60f33dda4db2cb1fcdcacf3410a8e05b3741f44a9b5998 AS alpine

ARG OS=linux
ARG ARCH=amd64

COPY --from=artifacts /bin/which_${OS}-${ARCH}-slim /bin/which

ARG VCS_REF
ARG VERSION
ARG CODEOWNERS

LABEL org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.title=go-which \
      org.opencontainers.image.authors=$CODEOWNERS \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.source="https://github.com/hairyhenderson/go-which"

ENTRYPOINT [ "/bin/which" ]

FROM scratch AS slim

ARG OS=linux
ARG ARCH=amd64

COPY --from=artifacts /bin/which_${OS}-${ARCH}-slim /which

ARG VCS_REF
ARG VERSION
ARG CODEOWNERS

LABEL org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.title=go-which \
      org.opencontainers.image.authors=$CODEOWNERS \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.source="https://github.com/hairyhenderson/go-which"

ENTRYPOINT [ "/which" ]
