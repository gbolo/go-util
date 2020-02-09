#
#  BUILD CONTAINER -------------------------------------------------------------
#

FROM gbolo/builder:debian as builder

COPY . /opt/gopath/src/cm
# Building
RUN   set -xe; \
      cd /opt/gopath/src/cm && ./build-bins.sh

#
#  FINAL BASE CONTAINER --------------------------------------------------------
#

FROM  gbolo/baseos:debian

# Copy in from builder
COPY  --from=builder /opt/gopath/src/cm/bin/* /
COPY  --from=builder /opt/gopath/src/cm/testdata /testdata

# Inherit gbolo/baseos entrypoint and pass it this argument
CMD  ["/agent"]
