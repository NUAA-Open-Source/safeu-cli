build:
	go build -o safeu
# On Mac OSX use this build binary file for Linux
linux-build:
	GOOS=linux GOARCH=amd64 go build -o safeu


