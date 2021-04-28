FROM iron/base
LABEL Description="This image is used to host the LT Text to IPA service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY data/text-to-ipa /app/
COPY config.yaml /app/

ENTRYPOINT ["./text-to-ipa"]
