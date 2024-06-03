all: build
build:
	bash ./build.sh
install:
	mkdir -p /etc/bestservers
	cp -n ./conf.ex.json /etc/bestservers/spy.conf
.PHONY: build