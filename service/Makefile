SELF := $(shell id -u)

build:
	@mkdir -p bin && go build -ldflags "-s -w" -gcflags "all=-trimpath=${PWD}" -o bin/craftythingshelper main.go

install:
	@install bin/craftythingshelper /usr/local/bin/ && \
		launchctl bootstrap gui/${SELF} ./com.voilaweb.fusion.craftythings.plist

uninstall:
	@launchctl bootout gui/${SELF} ./com.voilaweb.fusion.craftythings.plist

.PHONY: build install uninstall
