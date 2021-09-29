b:
	go build main.go
	mv main gpac
	cp gpac.gconf /etc/gpac.gconf
	install gpac /usr/bin
	mkdir -p /var/db/gpac/repo
	cp -rf repo /var/db/gpac/
	touch /var/db/gpac/pkgList
install: b