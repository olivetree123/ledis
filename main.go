package main

import (
	"encoding/json"
	"io"
	"ledis/commands"
	. "ledis/common"
	"ledis/handlers"
	"net"
)

func readFromConn(r io.Reader, responseChannel chan []byte, signalChannel chan bool) {
	decoder := json.NewDecoder(r)
	for {
		var command Command
		err := decoder.Decode(&command)
		if err == io.EOF {
			Logger.Infoln("End")
			break
		} else if err != nil {
			Logger.Errorln(err)
			break
		}
		Logger.Infoln("data = ", command)
		var result []byte
		if command.Code == commands.CreateTable {
			result = handlers.CreateTableHandler(command)
		} else if command.Code == commands.DeleteTable {
			result = handlers.DeleteTableHandler(command)
		} else if command.Code == commands.SetCache {
			result = handlers.SetCacheHandler(command)
		} else if command.Code == commands.GetCache {
			result = handlers.GetCacheHandler(command)
		}
		responseChannel <- result
	}
	signalChannel <- true
}

func handler(conn net.Conn) {
	defer conn.Close()
	var isClosed = false
	responseChannel := make(chan []byte)
	signalChannel := make(chan bool)
	go readFromConn(conn, responseChannel, signalChannel)
	for {
		select {
		case message := <-responseChannel:
			_, err := conn.Write(message)
			if err != nil {
				Logger.Errorln(err)
				break
			}
		case <-signalChannel:
			isClosed = true
			break
		}
		if isClosed {
			break
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		Logger.Errorln(err)
		return
	}
	Logger.Info("Start listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			Logger.Errorln(err)
			return
		}
		go handler(conn)
	}
}
