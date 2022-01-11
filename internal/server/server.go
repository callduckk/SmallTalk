package server

import (
	"fmt"
	"net"
	"sync"
)

type client struct {
	net.Conn
	name string
}

type Server struct {
	mu         sync.Mutex
	clientList []*client
	net.TCPListener
}

var Sv Server

func RunServer(endpoint string) error {
	if endpoint == "" {
		return createServerError(fmt.Errorf("endpoint cannot be empty"))
	}

	addr, err := net.ResolveTCPAddr("tcp", endpoint)

	if err != nil {
		return createServerError(fmt.Errorf("endpoint could not be parsed"))
	}

	l, err := net.ListenTCP("tcp", addr)

	if err != nil {
		return createServerError(fmt.Errorf("endpoint could not be listened"))
	}

	Sv := Server{}

	Sv.TCPListener = *l

	Sv.beginAcceptConnection(func(c *client) {
		go func() {
			for {
				err2 := c.receive(func(message string) {
					errs := Sv.broadcastMessage(message, c)
					for i := 0; i < len(errs); i++ {
						fmt.Println(errs[i])
					}
				})
				if err != nil {
					fmt.Println(err2)
				}
			}
		}()
	})

	fmt.Printf("Server started to run on: %s\n", endpoint)

	return nil
}

func (serv *Server) beginAcceptConnection(onAccept func(*client)) {
	go func() {
		by := make([]byte, 24)
		for {
			if len(serv.clientList) == 10 {
				break
			}
			acceptedCon, err := serv.AcceptTCP()
			if err != nil {
				continue
			}

			//TODO: make reading username a seperate function

			read, err := acceptedCon.Read(by)
			if err != nil {
				continue
			}
			name := string(by[:read])

			fmt.Printf("%s connected\n", name)

			cl := &client{acceptedCon, name}
			serv.mu.Lock()
			serv.clientList = append(serv.clientList, cl)
			serv.mu.Unlock()

			if onAccept != nil {
				onAccept(cl)
			}
		}
	}()
}

func (cl *client) receive(onReceive func(message string)) error {
	buffer := make([]byte, 0xFFFF)
	read, err := cl.Read(buffer)
	if err != nil {
		return fmt.Errorf("an error occurred while receiving a message from %s => %w", cl.name, err)
	}
	if read <= 0 {
		return fmt.Errorf("an error occurred while receiving a message from %s => %w", cl.name, fmt.Errorf("read len is zero"))
	}
	msg := string(buffer[:read])
	fmt.Printf("%s: %s\n", cl.name, msg)
	if onReceive != nil {
		onReceive(msg)
	}
	return nil
}

func (serv *Server) broadcastMessage(message string, sender *client) (errors []error) {
	for _, c := range serv.clientList {
		if c == sender {
			continue
		}
		sentMesgByte := []byte(sender.name + ": " + message)
		written, err := c.Write(sentMesgByte)
		if err != nil {
			errors = append(errors, fmt.Errorf("an error occurred while sending to %s => %w", c.name, err))
			continue
		}
		if written != len(sentMesgByte) {
			errors = append(errors, fmt.Errorf("could not send all of the message to %s", c.name))
		}
	}
	return
}

func createServerError(err error) error {
	return fmt.Errorf("an error occurred while creating Server: %w", err)
}
