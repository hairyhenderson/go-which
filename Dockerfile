FROM hairyhenderson/upx:3.96@sha256:68bd1a05f70602a58c52a7f4813f021cb2cefc41fe5449a4400361abf18c2b58 AS upx

FROM golang:1.27.3-alpine@sha256:fa3d596898d9f126af7c01d14ee7597bab3f39cccf392ef36b69fe77cd681f18 AS build

RUN apk add --no-cache \
      make \
      libgcc libstdc++ ucl \
      git

COPY --from=upx /usr/bin/upx /usr/bin/upx

RUN mkdir -p /go/src/github.com/hairyhenderson/go-which
WORKDIR /go/src/github.com/hairyhenderson/go-which
COPY . /go/src/github.com/hairyhenderson/go-which

ARG VCS_REF
ARG VERSION
ARG CODEOWNERS

RUN make build-x

FROM build AS compress

RUN make compress-all

FROM scratch AS artifacts

COPY --from=compress /go/src/github.com/hairyhenderson/go-which/bin/* /bin/

# gets around error 'Error response from daemon: No command specified'
CMD [ "/bin/which" ]

FROM scratch AS latest

ARG OS=linux
ARG ARCH=amd64

COPY --from=build /go/src/github.com/hairyhenderson/go-which/bin/which_${OS}-${ARCH} /which

ARG VCS_REF
ARG VERSION
ARG CODEOWNERS

LABEL org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.title=go-which \
      org.opencontainers.image.authors=$CODEOWNERS \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.source="https://github.com/hairyhenderson/go-which"

ENTRYPOINT [ "/which" ]

FROM alpine:3.24@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf AS alpine

ARG OS=linux
ARG ARCH=amd64

COPY --from=build /go/src/github.com/hairyhenderson/go-which/bin/which_${OS}-${ARCH} /bin/which

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

COPY --from=compress /go/src/github.com/hairyhenderson/go-which/bin/which_${OS}-${ARCH}-slim /bin/which

ARG VCS_REF
ARG VERSION
ARG CODEOWNERS

LABEL org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.title=go-which \
      org.opencontainers.image.authors=$CODEOWNERS \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.source="https://github.com/hairyhenderson/go-which"

ENTRYPOINT [ "/which" ]
