.PHONY: build generate test get install clean

all: get gen-client-secret build

build:
	go build

gen-client-secret:
	cd gdrive; go-bindata -pkg gdrive data

get:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u google.golang.org/api/drive/v3
	go get -u golang.org/x/oauth2/...

test:

install: get gen-client-secret
	go install

clean:
	-rm check-gdrive check-gdrive.exe bindata.go
