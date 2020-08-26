exe := sandstorm-rendezvous

all: $(exe)
pack: rendezvous.spk
dev: all
	spk dev

rendezvous.spk: all
	spk pack $@
sandstorm-rendezvous: $(wildcard *.go)
	CGO_ENABLED=0 go build
clean:
	rm -f $(exe) *.spk

.PHONY: all pack dev clean
