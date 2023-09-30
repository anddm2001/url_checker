.PHONY: build run

PATH_SOURCE_APP_FS=./cmd/main/main.go
PATH_APP_FS_WIN=./bin/checker/url_checker.exe

build:
	CGO_ENABLED=0 GOOS=windows go build -ldflags "-s -w" -o $(PATH_APP_FS_WIN) $(PATH_SOURCE_APP_FS)

run:
	(PATH_APP_FS_WIN)