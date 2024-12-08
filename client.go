package main

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var r = bufio.NewReader(os.Stdin)

func ConnectClient() error {
	//Create context with timeout for dial
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	//Construct url
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/",
	}

	//Dial to server (websocket handshake)
	c, _, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		return errors.New("conn error: server unreachable")
	}
	defer c.Close(websocket.StatusInternalError, "unexpected error, closing connection...")
	log.Printf("dial to %v successful", u.String())

	//Process signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	errc := make(chan error)

	//Incoming message listener
	go func() {
		var rec any
		for {
			if err := wsjson.Read(context.Background(), c, &rec); err != nil {
				errc <- err
				break
			}
			log.Printf("> received message: %v", rec)
		}
	}()

	//Message writer, gets message from stdin
	go func() {
		for {
			res, err := r.ReadString('\n')
			if err != nil {
				log.Println("failed to read message")
				continue
			}

			if err := wsjson.Write(context.Background(), c, res); err != nil {
				log.Printf("websocket client: %v", err)
				continue
			}

			log.Println("> message sent")
		}
	}()

	//Main loop
	for {
		//Waits for a signal
		select {
		case <-stop:
			log.Println("\nstopping client...")
			c.Close(websocket.StatusNormalClosure, "")
			return nil
		case err := <-errc:
			log.Printf("client error: %v", err)
			log.Println("stopping client...")
			c.Close(websocket.StatusNormalClosure, "")
			return nil
		}

	}
}
