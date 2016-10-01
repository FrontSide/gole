export GOPATH := $(HOME)/.go
export GOLE_PATH := $$(echo $$GOPATH | cut -d ":" -f 1)/src/gole
export PATH := $(PATH):$(GOPATH)/bin
export GOLE_GIT_PATH := $(pwd)

stop:
	kill -9 $$(cat server.pid) $$(cat client.pid) || true

clean: stop
	rm server.pid client.pid || true
	rm gole/gole || true
	cd gole && go clean

prepare:
	echo "Gopath is: $(GOPATH)"
	mkdir -p $$GOPATH/src || true
	echo "Gole Path is: $(GOLE_PATH)"
	rm -r $(GOLE_PATH) || true
	cp -r $$(pwd)/gole $$(echo $$GOPATH | cut -d ":" -f 1)/src/ 
	cd $(GOLE_PATH) && go get

test: prepare
	cd $(GOLE_PATH) && go test

build: clean prepare
	cd $(GOLE_PATH) && go build 
	cd $(GOLE_GIT_PATH)
	touch server.pid client.pid

start:
	gole & echo $$! >> server.pid 
	cd client 
	http-server -p 8080 & echo $$! >> client.pid
	cd ..

open: start
	open http://127.0.0.1:8080/client/game.html || xdg-open http://127.0.0.1:8080/client/game.html
