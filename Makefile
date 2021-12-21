.PHONY all: ./build/directChat

./build/directChat: ./build ./src/directChat.go
	go build -o ./build/directChat ./src/directChat.go
./build:
	mkdir build