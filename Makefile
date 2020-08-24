all: sandstorm-rendezvous
pack: rendezvous.spk
dev: all
	spk dev

rendezvous.spk: all
	spk pack $@
sandstorm-rendezvous: $(wildcard *.go)
	CGO_ENABLED=0 go build

.PHONY: all pack dev
