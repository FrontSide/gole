build:
	cd server && go clean && go build

start:
	server/server &
	open client/game.html || xdg-open client/game.html
