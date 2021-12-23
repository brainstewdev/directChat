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
	"encoding/gob"
)

var mode bool
var logging bool = true

var nickname string
var selected_color string

const SERVER = true
const CLIENT = false

var logger *log.Logger

var latest_messages []message

type message struct {
	Author	string
	Color string 
	Payload string
	Timestamp time.Time
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
func SetColor(rgb string) {
	//ESC[ 38;2;⟨r⟩;⟨g⟩;⟨b⟩ m Select RGB foreground color
	escape := "\033"
	fmt.Print(escape + "[38;2;" + rgb + "m")
}
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
func run_server(port string, hub connectionHub){ 
	fmt.Println("runno il server")
	listener := start_server(port)
	/*
	for i := 0; i<100; i++{
		time.Sleep(time.Second)
		mes := message{"lorenzo", strconv.Itoa(i), time.Now()}
		fmt.Println("mando...")
		rec <- mes
	}
	*/

	for {
		// accept the connections
		conn, err := (*listener).Accept()
			if logging{
				logger.Println("accepting connection...")
			}
		if err != nil {
			if logging {
				logger.Println("Error accepting connection:" ,err)
			}
		} else{
			// connection accepted succesfully, add the connection to the connections pool
			opened_conn := connection{conn,&hub} 
			hub.connections[&opened_conn] = true
			
			if logging {
				logger.Println("Connection accepted and added to connection pool");
			}
			// communicate to all client on the server that a new user has joined
			

			// ascolto sulla connessione, in caso arrivi qualcosa mando a tutti e stampo anche per debug
			go func(){				
				// gestisco la connessione
				dec := gob.NewDecoder(conn)
				for {
					if logging{
						logger.Println("Got message")
					}
					// prendo il messaggio che mi è stato inviato
					var msg message
					dec.Decode(&msg)
					if msg.Author == ""{
						// remove the connection from the connection pool
						delete(hub.connections, &opened_conn)
						conn.Close()
						break
					}
					
					fmt.Println("ricevuto", msg)
					
					if msg.Payload == "/quit"{
						// remove the user from the connection pool and close the connection.
						// then communicate this information to all the connected clients.
						delete(hub.connections, &opened_conn)
						conn.Close()
						msg = message{nickname, selected_color, "User " + msg.Author + " Left the chat.", time.Now()}	
					}else if msg.Payload == "/connectionrequest"{
						msg = message{nickname, selected_color, "User " + msg.Author + " Joined the chat.",  time.Now()}
					}
					// output the message to all client connected
					for c := range hub.connections {
						// c.realConn.Write(buf.Bytes())
						enc := gob.NewEncoder(c.realConn)
						enc.Encode(msg)
					}
				}
			}()
		}

	}
}
func run_client(address, port string){
	var waitgroup sync.WaitGroup
	if logging{
		logger.Println("Connecting to server with address", address, ":", port,"...")
	}
	conn, err := net.Dial("tcp", address+ ":"+ port)
	if err != nil {
		logger.Panic("Error connecting to server:", err)
	}else{
		waitgroup.Add(3)
		// defer conn.Close()
		messaggi_chan := make(chan message)

		// leggi nuovo messaggio
		go func(){
			// invia un messaggio che serve a comunicare al server che si è pronti per comunicare
			messaggi_chan<- message{nickname, selected_color, "/connectionrequest", time.Now()}

			for{
				scanner := bufio.NewScanner(os.Stdin)
				fmt.Print(">")
				scanner.Scan()
				payload := scanner.Text()
				
				msg := message{nickname, selected_color, payload, time.Now()}
				messaggi_chan <- msg
				if payload == "/quit"{
					break
				}
			}
			waitgroup.Done()
		}()
		// manda messaggio
		go func(){
			enc := gob.NewEncoder(conn)
			for {
				msg := <-messaggi_chan
				err := enc.Encode(msg)
				if err != nil{
					logger.Println("Error encoding struct", msg, ":", err)
				}
			}
			conn.Close()
			waitgroup.Done()
		}()
		// ricevi e stampa
		go func(){
			dec := gob.NewDecoder(conn)
			for {
				var msg message
				dec.Decode(&msg)

				if msg.Author != ""{
					// stampa tutti i messaggi partendo dal primo
					latest_messages = append(latest_messages, msg)
					ClearScreen()
					for _, v := range latest_messages{
						SetColor(v.Color)
						fmt.Println(v.Author, ":", v.Timestamp.Format(time.UnixDate), "\n>", v.Payload)
						SetColor("255;255;255");
					}
				}
			}
			waitgroup.Done()
		}()
	}
	waitgroup.Done()
}
func read_message(send chan message){
	var input string

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan(){
		input = scanner.Text()
		fmt.Println(":")
		send <- message{nickname, selected_color, input, time.Now()} 
	}
}

func print_message(rec <-chan message){
	for {
		select {
		case msg := <-rec:
			fmt.Println("new message:", msg)
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
	fmt.Println("starting..")
	if logging{
		logger.Output(2, "starting")
	}
	
	// default mode is client
	mode = CLIENT	
	// clean the screen..
	fmt.Print("\033[H\033[2J")
	// if args are passed and are exactly 3 then it is presumed that those args are correct and used instead of stdin
	var command_code string
	var ip_address string
	var port string
	if len(os.Args) == 3+1{
		command_code = os.Args[1]
		ip_address = os.Args[2]
		port = os.Args[3]	
	}else{
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
		command_code = command_parts[0]
		ip_address = command_parts[1]
		port = command_parts[2]
	}

	// creates the channels for the reiceved and to send message
	rec := make(chan message, 50)
	sen := make(chan message, 20)
	// create the waitgroup for the two goroutines
	var wg sync.WaitGroup
	wg.Add(3)
	// switch to run the commands inserted
	switch {
	case command_code == "c":
		fmt.Print("nickname:");
		fmt.Scan(&nickname)
		fmt.Print("Color(r;g;b):")
		fmt.Scan(&selected_color)
		fmt.Println("starting client mode")
		mode = CLIENT
		run_client(ip_address, port)
	case command_code == "s":
		nickname = "__SERVER__"
		selected_color = "255;255;255"
		fmt.Println("starting server mode")
		mode = SERVER
		var hub connectionHub
		hub.connections = make(map[*connection]bool)
		go func(){ 
			run_server(port, hub)
			wg.Done()
		}()
		go func(){
			print_message(rec)
			wg.Done()
		}()
		go func(){
			read_message(sen)
			wg.Done()
		}()
	}
	wg.Wait()
}
