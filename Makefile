VERSION := $(shell egrep "version *=" version.go | cut -d '"' -f 2)
GITHUB_VERSION := v$(VERSION)

all: build

fmt:
	go fmt $(shell find -name \*.go |xargs dirname|sort -u)

lint:
	golint $(shell find -name \*.go |xargs dirname|sort -u)

vet:
	go vet $(shell find -name \*.go |xargs dirname|sort -u)

build: fmt vet
	go build -asmflags -trimpath -o build/linux_amd64/letsencrypt-deploy

zip:
	mkdir dist
	zip -j dist/letsencrypt-deploy_$(GITHUB_VERSION)_linux_amd64.zip build/linux_amd64/letsencrypt-deploy

dist: clean build zip
	cd dist && sha512sum *.zip > letsencrypt-deploy_$(GITHUB_VERSION)_SHA512SUM.txt
	sed -e "s/<version>/$(VERSION)/g" -e "s/<checksum>/$(shell sha512sum dist/letsencrypt-deploy_$(GITHUB_VERSION)_linux_amd64.zip | cut -d " " -f 1)/g" terraform/deploy/variables.tf.json.template > terraform/deploy/variables.tf.json

clean:
	rm -rf build dist

sign:
	gpg --armor --sign --detach-sig dist/letsencrypt-deploy_$(GITHUB_VERSION)_linux_amd64.zip
	gpg --armor --sign --detach-sig dist/letsencrypt-deploy_$(GITHUB_VERSION)_SHA512SUM.txt

release:
	@echo "| File | Sign  | SHA512SUM |"
	@echo "|------|-------|-----------|"
	@echo "| [letsencrypt-deploy_$(GITHUB_VERSION)_linux_amd64.zip](../../releases/download/$(GITHUB_VERSION)/letsencrypt-deploy_$(GITHUB_VERSION)_linux_amd64.zip) | [letsencrypt-deploy_$(GITHUB_VERSION)_linux_amd64.zip.asc](../../releases/download/$(GITHUB_VERSION)/letsencrypt-deploy_$(GITHUB_VERSION)_linux_amd64.zip.asc) | $(shell sha512sum dist/letsencrypt-deploy_$(GITHUB_VERSION)_linux_amd64.zip | cut -d " " -f 1) |"

run:
	#./letsencrypt-deploy -email me@example.com -domain example.com,*.example.com -configFile config.json -o certificates/
	./run.sh
