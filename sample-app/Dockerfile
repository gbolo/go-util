#
#  BUILD CONTAINER -------------------------------------------------------------
#

FROM gbolo/builder:alpine as builder

COPY . /opt/gopath/src/sample-app
# Building
RUN   set -xe; \
      cd /opt/gopath/src/sample-app && make sample-app

#
#  FINAL BASE CONTAINER --------------------------------------------------------
#

FROM  gbolo/baseos:alpine

# Copy in from builder
COPY  --from=builder /opt/gopath/src/sample-app/bin/* /
COPY  --from=builder /opt/gopath/src/sample-app /testdata

# Inherit gbolo/baseos entrypoint and pass it this argument
CMD  ["/sample-app"]
