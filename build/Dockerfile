FROM iron/base

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

RUN addgroup app && adduser -S -G app app

RUN mkdir -p /app && chown -R app:app /app    

COPY data/text-to-ipa /app/
COPY config.yaml /app/
RUN chown app:app /app/*

USER app

ENTRYPOINT ["./text-to-ipa"]
