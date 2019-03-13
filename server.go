package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

const (
	protocol      = "tcp"
	nodeVersion   = 1
	commandLength = 12
)

var (
	nodeAddr   string
	miningAddr string
	knownNodes = []string{"localhost:3000"}
)

type verzion struct {
	Version    int    // blockchain version
	BestHeight int    // length of the node’s blockchain
	AddrFrom   string //address of the sender
}

type tx struct {
	AddFrom     string
	Transaction []byte
}

// StartServer starts a node.
func StartServer(nodeID, minerAddr string) {
	nodeAddr = fmt.Sprintf("localhost:%s", nodeID)
	miningAddr = minerAddr
	l, err := net.Listen(protocol, nodeAddr)
	if err != nil {
		log.Panic(err)
	}
	defer l.Close()

	bc := NewBlockChain(nodeID)

	if nodeAddr != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConn(conn, bc)
	}
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendTx(addr string, tnx *Transaction) {
	data := tx{
		AddFrom:     nodeAddr,
		Transaction: tnx.Serialize(),
	}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

func sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(
		verzion{
			Version:    nodeVersion,
			BestHeight: bestHeight,
			AddrFrom:   nodeAddr,
		},
	)

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte
	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func handleConn(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
	case "block":
	case "inv":
	case "getblocks":
	case "getdata":
	case "tx":
	case "version":
	default:
		fmt.Println("Unknown command!")
	}
	conn.Close()
}
