package client

import (
	"fmt"
	"net"
)

type Client struct {
	net.Conn
	name string
}

func NewClient(name string, endpoint string) (Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", endpoint)
	if err != nil {
		return Client{}, createClientError(err)
	}
	tcp, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return Client{}, createClientError(err)
	}
	_, err2 := tcp.Write([]byte(name))
	if err2 != nil {
		return Client{}, createClientError(err)
	}

	cl := Client{name: name, Conn: tcp}

	go func() {
		for {
			cl.Receive()
		}
	}()

	return cl, nil
}

func createClientError(err error) error {
	return fmt.Errorf("an error occurred while connecting: %w", err)
}

func (cl *Client) Receive() {
	b := make([]byte, 0xFFFF)
	read, err := cl.Read(b)
	if err != nil {

	}
	fmt.Printf("%s\n", string(b[:read]))
}

func (cl *Client) Send(message string) error {
	msgBytes := []byte(message)
	written, err := cl.Write(msgBytes)
	if err != nil {
		return sendMessageError(err)
	}
	if written != len(msgBytes) {
		return sendMessageError(fmt.Errorf("not all data could be send"))
	}
	return nil
}

func sendMessageError(err error) error {
	return fmt.Errorf("an error occurred while sending message: %w", err)
}
