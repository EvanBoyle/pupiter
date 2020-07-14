build: build-web pkger install

build-web: 
	cd ./notebook && yarn run build
pkger: 
	pkger
install:
	go install github.com/evanboyle/pupiter/cmd/pupiter