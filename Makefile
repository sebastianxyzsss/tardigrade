version = 0.0.3

.PHONY: build 

build:
	mkdir -p build
	CGO_ENABLED=0 go build -o build/tg main.go

clean:
	rm -rf build
	rm -rf tg/usr
	rm -rf tg/flatpak/build
	rm -rf buildfp

debiancp: build
	mkdir -p tg/usr/bin
	cp build/tg tg/usr/bin

debian: debiancp
	dpkg-deb --build tg
	mv tg.deb build/tg.$(version).deb

flatpak: build
	mkdir -p tg/flatpak/build
	cp build/tg tg/flatpak/build/tg
	flatpak-builder buildfp tg/flatpak/org.sebastian.xyzsss.Tardigrade.yml
