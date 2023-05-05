build:
	go build -ldflags "-s -w" -o bin/pm.exe cmd/main.go

create packet-2-2.0:
	pm create example/packet-2-2.0.json

create packet-1-1.00:
	pm create example/packet-1.00.json

create packet-1-1.10:
	pm create example/packet-1.10.json