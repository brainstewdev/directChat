.PHONY all: ./build/directChat

./build/directChat: ./build ./build/log ./src/directChat.go
	go build -o ./build/directChat ./src/directChat.go
./build:
	mkdir build
./build/log:
	mkdir ./build/log
clean:
	rm build -r