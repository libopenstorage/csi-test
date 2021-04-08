#! /bin/bash -e

fail() {
  echo "$1"
  exit 1
}

echo "Testing CSI Sanity against openstorage CSI Driver in current GOPATH"
set -x

# Copy GOPATH CSI Driver to csi-sanity for testing
cp -R vendor/github.com/libopenstorage/openstorage/csi vendor/github.com/libopenstorage/openstorage/csi-saved || fail "failed save CSI driver copy from local repo"
cp -R $GOPATH/src/github.com/libopenstorage/openstorage/csi vendor/github.com/libopenstorage/openstorage || fail "failed to copy CSI driver from GOPATH"

# GOPATH version of CSI Driver to csi-sanity and test
go test -v