stop:
	kill -9 $$(cat server.pid) $$(cat client.pid) || true

clean: stop
	rm server.pid client.pid || true
	rm gole/gole || true
	cd gole && go clean

test: 
	cd gole && go test

build: clean
	cd gole && go build 
	cd ..
	touch server.pid client.pid

start:
	gole/gole & echo $$! >> server.pid 
	cd client 
	http-server -p 8080 & echo $$! >> client.pid
	cd ..

open: start
	open http://127.0.0.1:8080/client/game.html || xdg-open http://127.0.0.1:8080/client/game.html
