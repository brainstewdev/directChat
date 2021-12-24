# directChat
A simple Chat software from within the Shell
## Usage
### Build
The client and the server are within the same executable. in order to build the binary you just need to use the make tool.
it will build the binary inside the build directory and will create the log folder.
### Server
The server allows multiple users to connect and chat togheter. If you want to chat with your friend and no server are available you need to start your own. to do that you need to build the program, run it and then input the s command, followed by any string (it is not used) followed by the port number.
this will start the server.
note: it may be neccessary to configure your network if you want to use the server with friends which are not inside your local network.
### Client
The Client can connect to a server and chat with other clients. you need to set a username and a color for how your text display before being able to connect to a server.
to start the software in client mode you just need to insert the c command (which stands for connect) followed by the ip address of the server followed by the port of the server.
