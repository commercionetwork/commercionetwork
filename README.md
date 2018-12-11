# Build the app
## Update dependencies to match the constraints and overrides above
dep ensure -update -v

## Install the app into your $GOBIN
make install

## Now you should be able to run the following commands:
nsd help
nscli help