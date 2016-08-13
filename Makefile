stop:
	kill -9 $$(cat server.pid) $$(cat client.pid) || true

clean: stop
	rm server.pid client.pid || true
	rm server/server || true
	cd server && go clean

build: clean
	cd server && go build 
	cd ..
	touch server.pid client.pid

start:
	server/server & echo $$! >> server.pid 
	cd client 
	http-server -p 8080 & echo $$! >> client.pid
	cd ..

open: start
	open http://127.0.0.1:8080/game.html || xdg-open http://127.0.0.1:8080/game.html
