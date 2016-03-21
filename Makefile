
fsbuilder::
	go build

output:: fsbuilder
	./fsbuilder
	
beta:: output
	rsync -av output/. gigo.com:/var/www/beta.test-ipv6.com/. --exclude site --delete
	ssh gigo.com "cd /var/www/beta.test-ipv6.com/ && ls -lt | grep -v orig"

setup:
	rsync -av gigo.com:falling-sky/source/templates/. templates
	rsync -av gigo.com:falling-sky/source/translations/. translations
	rsync -av gigo.com:falling-sky/source/images/. images
	rsync -av gigo.com:falling-sky/source/transparent/. transparent

debug::
	godebug run \
		-instrument github.com/falling-sky/builder/config \
		-instrument github.com/falling-sky/builder/fileutil \
		-instrument github.com/falling-sky/builder/job \
		-instrument github.com/falling-sky/builder/po \
		builder.go
