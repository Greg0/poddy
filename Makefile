build-osx:
	env GOOS=darwin GOARCH=amd64 go build -o build/osx/poddy poddy.go

build-win:
	env GOOS=windows GOARCH=amd64 go build -o build/win/poddy.exe poddy.go

build-lin:
	env GOOS=linux GOARCH=amd64 go build -o build/lin/poddy poddy.go

build: build-lin build-osx build-win