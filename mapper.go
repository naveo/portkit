package main

import (
	"context"
	"io"
	"log"
	"net"
)

func PortMapper(id, port string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("unprivileged bind ports lower than 1024 require privilege escalation")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, eventChannel, errChannel := DockerEventListener()

	listen, err := net.Listen("tcp", config.TargetHost+":"+port)
	if err != nil {
		log.Printf("close listener server for port %v with error %v", port, err)
	}
	defer listen.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				hostConn, err := listen.Accept()
				if err != nil {
					log.Printf("stop listening on port %v", port)
				}
				go handleConnection(ctx, hostConn, port)
			}

		}

	}()

	for {
		select {
		case err := <-errChannel:
			log.Fatalf("error reading events channel %v:", err)
		case containerEvent := <-eventChannel:
			if containerEvent.Status == "die" && containerEvent.ID == id {
				log.Printf("end port forward for cid %v port %v", id, port)
				cancel()
				return
			}
		}
	}
}

func handleConnection(ctx context.Context, hostConn net.Conn, port string) {
	guestConn, err := net.Dial("tcp", config.NaveoSocket+":"+port)
	if err != nil {
		log.Printf("port mapping closed for port %v error %v", port, err)
	}

	if guestConn != nil {
		select {
		case <-ctx.Done():
			return
		default:
			go handleTraffic(ctx, hostConn, guestConn)
			go handleTraffic(ctx, guestConn, hostConn)
		}
	}
}

func handleTraffic(ctx context.Context, source, destination net.Conn) {
	select {
	case <-ctx.Done():
		return
	default:
		defer source.Close()
		defer destination.Close()
		io.Copy(source, destination)
	}
}
