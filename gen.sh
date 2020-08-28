#!/usr/bin/env sh

set -euo pipefail
# To run this, you need to set the variables to point to the paths
# containing capnp files from go.sandstorm and go-capnproto2, respectively.
capnp compile -ogo \
	-I $SANDSTORM_CAPNP \
	-I $GO_CAPNP \
	*.capnp
