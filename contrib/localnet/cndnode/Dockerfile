FROM ubuntu:18.04

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get -y install curl jq file tzdata

VOLUME /cnd
WORKDIR /cnd
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM
ENV TZ America/New_York

COPY wrapper.sh /usr/bin/wrapper.sh
