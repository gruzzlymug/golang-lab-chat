package main

import (
    "bufio"
	"log"
	"net"
)

type ChatRoom struct {
	// TODO: populate this
    users map[string]*ChatUser
    incoming chan string
    joins chan *ChatUser
    disconnects chan string
}

// NewChatRoom will create a chatroom
func NewChatRoom() *ChatRoom {
	// TODO: initialize struct members
	return &ChatRoom{
        users: make(map[string]*ChatUser),
        incoming: make(chan string),
        joins: make(chan *ChatUser),
        disconnects: make(chan string),
    }
}

func (cr *ChatRoom) ListenForMessages() {
    go func() {
        for {
            select {
            case msg := <-cr.incoming:
                cr.Broadcast(msg)
            case user := <-cr.joins:
                cr.users[user.username] = user
                cr.Broadcast("*** " + user.username + " just joined the chatroom")
            }
        }
    }()
}

func (cr *ChatRoom) Logout(username string) {}
func (cr *ChatRoom) Join(conn net.Conn) {
    user := NewChatUser(conn)
    err := user.Login(cr)
    if err != nil {
        log.Fatal("Could not login user")
    }
    cr.joins <- user
}

func (cr *ChatRoom) Broadcast(msg string) {
    for _, user := range cr.users {
        user.Send(msg)
    }
}

type ChatUser struct {
	// TODO: populate this
    conn net.Conn
    disconnect bool
    username string
    outgoing chan string
    reader *bufio.Reader
    writer *bufio.Writer
}

func NewChatUser(conn net.Conn) *ChatUser {
	// TODO: initialize chat user
	return &ChatUser{
        conn: conn,
        disconnect: false,
        username: "",
        outgoing: make(chan string),
        reader: bufio.NewReader(conn),
        writer: bufio.NewWriter(conn),
    }
}

func (cu *ChatUser) ReadIncomingMessages(chatroom *ChatRoom) {
	// TODO: read incoming messages in a loop
    go func() {
        for {
            line, _ := cu.ReadLine()
            if line != "" {
                chatroom.incoming <- ("[" + cu.username + "] " + line)
            }
        }
    }()
}

func (cu *ChatUser) WriteOutgoingMessages(chatroom *ChatRoom) {
	// TODO: wait for outgoing messages in a loop, and write them
    go func() {
        for {
            data := <-cu.outgoing
            data = data + "\n"
            cu.WriteString(data)
        }
    }()
}

func (cu *ChatUser) Login(chatroom *ChatRoom) error {
	// TODO: login the user
    cu.WriteString("Welcome the the chat server\n")
    cu.WriteString("Please enter your username: ")

    username_bytes, _, _ := cu.reader.ReadLine()
    cu.username = string(username_bytes)

    cu.WriteOutgoingMessages(chatroom)
    cu.ReadIncomingMessages(chatroom)

    cu.WriteString("Welcome " + cu.username + "\n")
	return nil
}

func (cu *ChatUser) ReadLine() (string, error) {
	// TODO: read a line from the socket
    bytes, _, err := cu.reader.ReadLine()
    str := string(bytes)
	return str, err
}

func (cu *ChatUser) WriteString(msg string) error {
	// TODO: write a line from the socket
    cu.writer.WriteString(msg)
    cu.writer.Flush()
	return nil
}

func (cu *ChatUser) Send(msg string) {
	// TODO: put a message on the outgoing messages queue
    cu.outgoing <- msg
}

func (cu *ChatUser) Close() {
	// TODO: close the socket
}

//
// main will create a socket, bind to port 6677,
// and loop while waiting for connections.
//
// When it receives a connection it will pass it to
// `chatroom.Join()`.
//
func main() {
	log.Println("Chat server starting!")
    room := NewChatRoom()

	// TODO add other logic
    listener, err := net.Listen("tcp", ":6677")
    if err != nil {
        log.Fatal("Unable to bind to 6677", err)
    }

    room.ListenForMessages()

    for {
        conn, err := listener.Accept()
        if err != nil {
            // handle error
        }

        go room.Join(conn)
        
        log.Println("Connection joined: ", conn.RemoteAddr())
    }
}
