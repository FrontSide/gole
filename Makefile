clean:
	rm server/server || true
	cd server && go clean
build: clean
	cd server && go build

start:
	server/server &
	cd client && http-server -p 8080 || true &

open: start
	open http://127.0.0.1:8080/game.html || xdg-open http://127.0.0.1:8080/game.html
