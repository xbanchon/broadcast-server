package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var broadcast = make(chan any)

func StartServer(port int) error {
	addr := fmt.Sprintf("localhost:%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("listening on ws://%v", l.Addr())

	s := &http.Server{
		Handler: Broadcaster{
			broadcast: make(chan string),
			clients:   make(map[*websocket.Conn]bool),
			// register:   make(chan *websocket.Conn),
			// unregister: make(chan *websocket.Conn),
		},
	}

	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case <-stop:
		log.Println("\nserver stopped.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Shutdown(ctx)
}

type Broadcaster struct {
	clients   map[*websocket.Conn]bool
	broadcast chan string
}

func (b Broadcaster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Upgrades connection (http -> ws)
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer c.CloseNow()

	if c.Subprotocol() != "echo" {
		c.Close(websocket.StatusPolicyViolation, "client must speak the echo protocol")
		return
	}

	//Adds connection to clients list if not in list
	_, ok := b.clients[c]
	if !ok {
		b.clients[c] = true
	}

	go func() {
		for {
			message, err := readMessage(c)
			if err != nil {
				log.Println("error reading message")
				log.Println(err)
				delete(b.clients, c)
				log.Print("closing client connection...\n\n")
				c.Close(websocket.StatusInternalError, "")
				break
			}

			broadcast <- message
		}
	}()

	for msg := range broadcast {
		for client := range b.clients {
			if err := echoMessage(client, msg); err != nil {
				log.Println("error echoing message")
				continue
			}
		}
	}
}

func readMessage(c *websocket.Conn) (any, error) {
	var msg any
	if err := wsjson.Read(context.Background(), c, &msg); err != nil {
		return "", err
	}
	return msg, nil
}

func echoMessage(c *websocket.Conn, msg any) error {
	if err := wsjson.Write(context.Background(), c, msg); err != nil {
		return err
	}
	return nil
}
