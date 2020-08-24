# FIXME: this will only work on my(zenhack) machine, due to the hard-coded
# paths we're passing to -I.
capnp compile -ogo \
	-I ~/src/pub/go.sandstorm/capnp/ \
	-I ~/src/foreign/go-capnproto2/std/ \
	*.capnp
