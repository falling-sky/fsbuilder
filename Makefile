
builder::
	go build

output:: builder
	./builder
	
beta:: output
	rsync -av output/. gigo.com:/var/www/beta.test-ipv6.com/. --exclude site --delete
	ssh gigo.com "cd /var/www/beta.test-ipv6.com/ && ls -lt | grep -v orig"

debug::
	godebug run \
		-instrument github.com/falling-sky/builder/config \
		-instrument github.com/falling-sky/builder/fileutil \
		-instrument github.com/falling-sky/builder/job \
		-instrument github.com/falling-sky/builder/po \
		builder.go
