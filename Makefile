b:
	go build main.go
	mv main gpac
	cp gpac.gconf /etc/gpac.gconf
	install gpac /usr/bin
	mkdir -p /var/db/gpac/repo
	cp -rf grepo /var/db/gpac/
install:
	go build main.go
	mv main gpac
	cp gpac.gconf /etc/gpac.gconf
	install gpac /usr/bin
	mkdir -p /var/db/gpac/repo
	cp -rf grepo /var/db/gpac/
s:
	@echo "put the path to repo in gpac.conf (it ends with a /)"


