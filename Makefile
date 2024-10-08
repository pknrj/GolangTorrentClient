build : 
	go build -o bin/GolangTorrentClient

run : build
	./bin/GolangTorrentClient