SHELL:=/bin/bash

build:
	go build -o gasbot

run: build
	./gasbot
