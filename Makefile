VERSION := v$(shell egrep "version *=" version.go | cut -d '"' -f 2)

all: build

fmt:
	go fmt $(shell find -name \*.go |xargs dirname|sort -u)

vet:
	go vet $(shell find -name \*.go |xargs dirname|sort -u)

build: fmt vet
	go build -asmflags -trimpath -o build/linux_amd64/letsencrypt-deploy

dist: clean build
	mkdir dist
	zip -j dist/letsencrypt-deploy_$(VERSION)_linux_amd64.zip build/linux_amd64/letsencrypt-deploy
	cd dist && sha512sum *.zip > letsencrypt-deploy_$(VERSION)_SHA512SUM.txt

clean:
	rm -rf build dist

sign:
	gpg --armor --sign --detach-sig dist/letsencrypt-deploy_$(VERSION)_linux_amd64.zip

release:
	@echo "| File | Sign  | SHA512SUM |"
	@echo "|------|-------|-----------|"
	@echo "| [letsencrypt-deploy_$(VERSION)_linux_amd64.zip](../../releases/download/$(VERSION)/letsencrypt-deploy_$(VERSION)_linux_amd64.zip) | [letsencrypt-deploy_$(VERSION)_linux_amd64.zip.asc](../../releases/download/$(VERSION)/letsencrypt-deploy_$(VERSION)_linux_amd64.zip.asc) | $(shell sha512sum dist/letsencrypt-deploy_$(VERSION)_linux_amd64.zip | cut -d " " -f 1) |"

run:
	#echo <secure_client_passphrase> > /tmp/deploy.passphrase
	#./letsencrypt-deploy -email me@example.com -domain example.com,*.example.com -passphraseFile /tmp/deploy.passphrase -o certificates/
	#rm -f /tmp/deploy.passphrase
	./run.sh
