b:
	go build main.go
	mv main gpac
	cp gpac.gconf /etc/gpac.gconf
	install gpac /usr/bin
	mkdir -p /var/db/gpac/repo
	cp -rf repo /var/db/gpac/
	touch /var/db/gpac/pkgList
install: b

runNeofetch:
	rm -f /usr/bin/neofetch
	/usr/bin/gpac b neofetch

runPfetch:
	rm -f /usr/bin/pfetch
	/usr/bin/gpac b legacy/pfetch

checkPfetch:
	rm /usr/bin/pfetch

checkNeofetch:
	rm /usr/bin/neofetch