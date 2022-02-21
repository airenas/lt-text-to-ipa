#########################################################################################
## docker will invoke this file from ../ dir in order to access code
#########################################################################################
FROM golang:1.17-alpine as builder

ARG BUILD_VERSION=0.1

WORKDIR /go/src/
ENV CGO_ENABLED=0

COPY . /go/src
RUN CGO_ENABLED=0 go build -o /go/bin/text-to-ipa -ldflags "-X main.version=$BUILD_VERSION" cmd/text-to-ipa/main.go
#########################################################################################
FROM iron/base:3.14.0

ARG BUILD_VERSION=0.1

LABEL org.opencontainers.image.version=$BUILD_VERSION \
      org.opencontainers.image.authors="airenass@gmail.com" \
      name="sampa-to-ipa" \
      maintainer="airenass@gmail.com" \
      vendor="airenass@gmail.com" \
      version=$BUILD_VERSION \
      release=$BUILD_VERSION \
      summary="This image is used to host the LT Text to IPA service" \
      description="LT Text to IPA service runtime" 

WORKDIR /app
EXPOSE 8000
ENV UID=999

RUN addgroup app && adduser -S -G app app --uid $UID app
RUN mkdir -p /app && chown -R app:app /app   


COPY --from=builder /go/bin/text-to-ipa /app/
COPY build/config.yaml /app/

RUN chown app:app /app/*

USER app

ENTRYPOINT ["./text-to-ipa"]
