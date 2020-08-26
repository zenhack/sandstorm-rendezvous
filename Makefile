exe := sandstorm-rendezvous
cleanfiles := \
	static/vnc-bundle.min.js \
	ui/vnc-bundle.js \
	$(exe)

all: $(exe) static/vnc-bundle.min.js
pack: rendezvous.spk
dev: all
	spk dev

rendezvous.spk: all
	spk pack $@
sandstorm-rendezvous: $(wildcard *.go)
	CGO_ENABLED=0 go build
ui/vnc-bundle.js: ui/index.js $(wildcard ui/package*.json)
	cd ui && npx rollup --config
static/vnc-bundle.min.js: ui/vnc-bundle.js
	(cd ui && npx uglifyjs --compress --mangle) < $< > $@
clean:
	rm -f $(cleanfiles) *.spk

.PHONY: all pack dev clean
