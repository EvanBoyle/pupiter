build: build-web pkger install

build-web: 
	cd ./notebook && yarn run build
pkger:
	pkger
	pkger -o cmd/pupiter
install:
	go install github.com/evanboyle/pupiter/cmd/pupiter