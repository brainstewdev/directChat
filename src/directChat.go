package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"net"
	"log"
	"time"
	"sync"
	"strconv"
)

var mode bool
var logging bool = true

const SERVER = true
const CLIENT = false

var logger *log.Logger

type message struct {
	author	string 
	payload string
	timestamp time.Time
}

type connectionHub struct {
	connections map[*connection]bool
}

type connection struct {
	realConn net.Conn
	hub *connectionHub
}

func start_server(port string) *net.Listener {
	ln, err := net.Listen("tcp", ":"+ port)
	if err != nil {
		if logging{
			logger.Fatal("ERROR starting server")
		}
	}else{
		if logging{
			logger.Output(2, "Started the Server on port " + port)
		}
	}
	return &ln
}
func start_client(){
	
}
// function which handle a connection
// it listen to client sent messages and then send the messages to all the clients
// connected
func handle_connection(conn Connection){

}
func run_server(port string, rec chan message, sen <-chan message){
	fmt.Println("runno il server")
	listener := start_server(port)
	for i := 0; i<100; i++{
		time.Sleep(time.Second)
		mes := message{"lorenzo", strconv.Itoa(i), time.Now()}
		fmt.Println("mando...")
		rec <- mes
	}
	close(rec)
	_ = listener
}
func run_client(){

}
func read_message(send chan message){
	var input string

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan(){
		input = scanner.Text()
		fmt.Println(":")
		send <- message{"lorenzo", input, time.Now()} 
	}
}

func print_message(rec <-chan message, sen <- chan message){
	for {
		select {
		case msg := <-rec:
				fmt.Println(msg)
		case msg := <-sen:
				fmt.Println(msg)
		}
	}
}

func main() {
	
	// open the log file
	out_log, err := os.Create("./log" + string(os.PathSeparator) + strings.ReplaceAll(time.Now().Format(time.UnixDate), " ", "_") + ".log")
	if err != nil{
		fmt.Println("ERROR opening log file, will not be logging on file")
		fmt.Println("file opening error:", err)
		logging = false
	}else{
		logger = log.New(out_log, "directChat ", log.Ldate | log.Ltime)
	}
	// default mode is client
	mode = CLIENT

	fmt.Println("starting..")
	if logging{
		logger.Output(2, "starting")
	}
	// clean the screen..

	fmt.Print(">")
	var command string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	command = scanner.Text()

	command = strings.Split(command, "\n")[0] // removes the \n from the command
	command_parts := strings.Split(command, " ")

	// every command follow the convention:
	/*
		command code (c for connect s to start a server)
		ip address (not used in caso of s command)
		port
	*/
	command_code := command_parts[0]
	ip_address := command_parts[1]
	port := command_parts[2]

	// creates the channels for the reiceved and to send message
	rec := make(chan message, 50)
	sen := make(chan message, 20)
	// create the waitgroup for the two goroutines
	var wg sync.WaitGroup
	wg.Add(3)
	// switch to run the commands inserted
	switch {
	case command_code == "c":
		fmt.Println("starting client mode")
		mode = CLIENT
	case command_code == "s":
		fmt.Println("starting server mode")
		mode = SERVER
		go func(){ 
			run_server(port, rec, sen)
			wg.Done()
		}()
		go func(){
			print_message(rec, sen)
			wg.Done()
		}()
		go func(){
			read_message(sen)
			wg.Done()
		}()
	}
	_ = port
	_ = ip_address
	wg.Wait()
}
