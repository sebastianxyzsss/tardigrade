version = 0.0.2

.PHONY: build 

build:
	go build -o build/tg main.go

clean:
	rm -rf build
	rm -rf tg/usr

debiancp: build
	mkdir -p tg/usr/bin
	cp build/tg tg/usr/bin

debian: debiancp
	dpkg-deb --build tg
	mv tg.deb build/tg.$(version).deb
