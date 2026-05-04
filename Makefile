GO    		:= go
glide    	:= glide
release    	:= ./release.sh
.PHONY: all clean appmanifest

all: clean buildall

clean: 
	rm -rf ./build/*

buildall: appmanifest

appmanifest: 
	rm -rf ./build/appmanifest ./build/appmanifest-darwin-amd64 ./build/appmanifest-darwin-arm64
	@echo ">> building appmanifest"
	cd ./appmanifest && $(release)
