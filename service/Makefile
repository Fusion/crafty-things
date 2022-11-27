SELF := $(shell id -u)

build:
	@mkdir -p bin && go build -ldflags "-s -w" -gcflags "all=-trimpath=${PWD}" -o bin/craftythingshelper main.go

install:
	@install bin/craftythingshelper /usr/local/bin/ && \
		launchctl bootstrap gui/${SELF} ./com.voilaweb.fusion.craftythings.plist

uninstall:
	@launchctl bootout gui/${SELF} ./com.voilaweb.fusion.craftythings.plist

packageapp:
	@mkdir -p dist/work \
	&& rm -rf dist/work/craftythingshelper.app  \
	&& cp -r packaging/* dist/work/  \
	&& mkdir -p dist/work/craftythingshelper.app/Contents/MacOS  \
	&& cp bin/craftythingshelper dist/work/craftythingshelper.app/Contents/MacOS/ \
	&& rm -f /tmp/tmp.dmg ./dist/CraftyThingsHelper.dmg \
	&& hdiutil create /tmp/tmp.dmg -ov -volname "CraftyThingsHelper" -fs HFS+ -srcfolder "./dist/work" \
	&& hdiutil convert /tmp/tmp.dmg -format UDZO -o ./dist/CraftyThingsHelper.dmg

.PHONY: build install uninstall packageapp
