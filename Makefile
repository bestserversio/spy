all: build
build:
	bash ./build.sh
install:
	mkdir -p /etc/bestservers
	cp -n ./conf.ex.json /etc/bestservers/spy.json
	cp -f ./build/spy /usr/bin/spy
.PHONY: build