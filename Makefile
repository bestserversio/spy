all: install
install:
	mkdir -p /etc/bestservers
	cp -n ./conf.ex.json /etc/bestservers/spy.conf