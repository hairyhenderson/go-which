FROM alpine:3.8 AS upx
RUN apk add --no-cache upx=3.94-r0

FROM golang:1.13.5-alpine AS build

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

FROM alpine:3.11@sha256:3983cc12fb9dc20a009340149e382a18de6a8261b0ac0e8f5fcdf11f8dd5937e AS alpine

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
