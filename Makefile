export GOPATH := $(GOPATH):$(PWD)/server
export SERVER_SRC := $(PWD)/server/src
export BUILD_PATH := $(PWD)

stop:
	kill -9 $$(cat server.pid) $$(cat client.pid) || true

clean: stop
	rm server.pid client.pid || true
	rm gole || true
	cd $(SERVER_SRC) && go clean

prepare-server:
	echo "Goroot is: $(GOROOT)"
	echo "Gopath is: $(GOPATH)"

prepare-client:
	cd client && npm install

prepare: prepare-server prepare-client

test: prepare-server
	cd $(SERVER_SRC) && go test

build: clean prepare
	cd $(SERVER_SRC) && go build -o gole
	mv $(SERVER_SRC)/gole $(BUILD_PATH)
	touch server.pid client.pid

start-server:
	./gole & echo $$! >> server.pid 

start-client:	
	cd client && npm start & echo $$! >> ../client.pid
	cd ..

start: start-server start-client
