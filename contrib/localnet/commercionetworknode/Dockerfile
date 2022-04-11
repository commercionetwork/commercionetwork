#FROM golang:latest 
FROM commercionetwork/libraries

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN make build

VOLUME /commercionetwork
WORKDIR /commercionetwork
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM
ENV TZ America/New_York

COPY ./contrib/localnet/commercionetworknode/wrapper.sh /usr/bin/wrapper.sh