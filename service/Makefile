SELF := $(shell id -u)

buildamd64:
	@GOARCH=amd64 go build -ldflags "-s -w" -gcflags "all=-trimpath=${PWD}" -o bin/craftythingshelper-amd64 main.go

buildarm64:
	@GOARCH=arm64 go build -ldflags "-s -w" -gcflags "all=-trimpath=${PWD}" -o bin/craftythingshelper-arm64 main.go

build:
	@mkdir -p bin \
		&& make buildamd64 \
		&& make buildarm64 \
		&& lipo -create -output bin/craftythingshelper bin/craftythingshelper-amd64 bin/craftythingshelper-arm64 \
		&& rm bin/craftythingshelper-*

install:
	@install bin/craftythingshelper /usr/local/bin/ && \
		launchctl bootstrap gui/${SELF} ./com.voilaweb.fusion.craftythings.plist

uninstall:
	@launchctl bootout gui/${SELF} ./com.voilaweb.fusion.craftythings.plist

packageapp:
	@mkdir -p dist/work \
	&& rm -rf dist/work/craftythingshelper.app  \
	&& cp -r packaging/* dist/work/  \
	&& ln -sf /Applications dist/work/Applications \
	&& mkdir -p dist/work/craftythingshelper.app/Contents/MacOS  \
	&& cp bin/craftythingshelper dist/work/craftythingshelper.app/Contents/MacOS/ \
	&& rm -f /tmp/tmp.dmg ./dist/CraftyThingsHelper.dmg \
	&& hdiutil create /tmp/tmp.dmg -ov -volname "CraftyThingsHelper" -fs HFS+ -srcfolder "./dist/work" \
	&& hdiutil convert /tmp/tmp.dmg -format UDZO -o ./dist/CraftyThingsHelper.dmg

.PHONY: build buildamd64 buildarm64 install uninstall packageapp
