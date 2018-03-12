all:
	go install

release:
	export GOOS=darwin GOARCH=amd64; go build -o gemacs_osx
	export GOOS=windows GOARCH=amd64; go build -o gemacs_windows.exe
	export GOOS=linux GOARCH=amd64; go build -o gemacs_linux
