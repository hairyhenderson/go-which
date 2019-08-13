FROM golang:1.12.8-alpine@sha256:057690933dda0237cef8315329d3f35c7d135c6fcffa072883a84e23988cab2e AS build

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
