package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/urfave/cli/v3"
	"golang.org/x/net/proxy"
)

func main() {
	cmd := &cli.Command{
		Name: "conn",
		Commands: []*cli.Command{
			{
				Name:  "c",
				Usage: "Establish a TCP connection to a specified address",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "s5",
						Usage: "Specify a SOCKS5 proxy address for the connection",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					if len(c.Args().Slice()) > 1 {
						return errors.New("The 'c' command only accepts one parameter")
					}
					if len(c.Args().Slice()) < 1 {
						return errors.New("The 'c' command expects at least one parameter")
					}

					connAddr := c.Args().Get(0)
					socks5Addr := c.String("s5")

					return connect(connAddr, socks5Addr)
				},
			},
			{
				Name:  "l",
				Usage: "Start a TCP server and listen on the specified port",
				Action: func(ctx context.Context, c *cli.Command) error {
					if len(c.Args().Slice()) > 1 {
						return errors.New("The 'l' command only accepts one parameter")
					}
					if len(c.Args().Slice()) < 1 {
						return errors.New("The 'l' command expects at least one parameter")
					}
					port := c.Args().Get(0)
					return listen(port)
				},
			},
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func connect(connAddr string, socks5Addr string) error {
	var conn net.Conn
	var err error

	if socks5Addr != "" {
		dialer, err := proxy.SOCKS5("tcp", socks5Addr, nil, proxy.Direct)
		if err != nil {
			return fmt.Errorf("SOCKS5 proxy error: %v", err)
		}

		conn, err = dialer.Dial("tcp", connAddr)
		if err != nil {
			return fmt.Errorf("Connection error: %v", err)
		}
	} else {
		conn, err = net.Dial("tcp", connAddr)
		if err != nil {
			return fmt.Errorf("Connection error: %v", err)
		}
	}

	defer conn.Close()

	return handleConn(conn)
}

func listen(port string) error {
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return err
	}

	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Println("Error accepting connection: ", err)
	}
	fmt.Println("New connection: ", conn.RemoteAddr().String())

	return handleConn(conn)
}

func handleConn(conn net.Conn) error {

	connClose := make(chan struct{})
	stdinCh := make(chan string)

	go func(connClose chan<- struct{}) {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Connection closed.")
				connClose <- struct{}{}
				break
			}
			fmt.Print(string(buf[:n]))
		}
	}(connClose)

	go func(stdinCh chan<- string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				break
			}
			stdinCh <- input
		}
	}(stdinCh)

	for {
		select {
		case <-connClose:
			return nil
		case input := <-stdinCh:
			_, err := conn.Write([]byte(input))
			if err != nil {
				fmt.Println("Error writing to connection:", err)
				return err
			}
		}
	}
}
