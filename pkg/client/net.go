package client

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"time"
)

const (
	maxDatagramSize = 8192
)

func parseReply(buffer []byte, key string) (reply *Reply, err error) {

	if err := json.Unmarshal(buffer, &reply); err != nil {
		return reply, err
	}

	reply.Pack, err = reply.DecryptPack(key)

	return
}

func newMessage(req RequestInterface, key string, controlRequest bool) (message []byte, err error) {

	encrypted, err := req.GetEncrypted([]byte(key))
	if err != nil {
		return
	}

	i := uint(1)
	if controlRequest {
		i = 0
	}

	message, err = json.Marshal(Message{
		ClientID: "app",
		I:        i,
		Pack:     base64.StdEncoding.EncodeToString(encrypted),
		MAC:      req.GetMAC(),
		Type:     "pack",
		UID:      0,
	})

	return
}

func listen(address string) (conn *net.UDPConn, err error) {

	// Parse the string address
	listenAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return
	}

	// Open up a connection
	conn, err = net.ListenUDP("udp", listenAddr)
	return
}

func send(address string, message []byte) (resp []byte, err error) {

	conn, err := net.Dial("udp", address)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.Write(message)
	conn.SetReadDeadline(time.Now().Add(time.Second * 1))

	buffer := make([]byte, maxDatagramSize)
	numBytes, err := conn.Read(buffer)

	resp = buffer[0:numBytes]

	return
}

func readForAWhile(conn *net.UDPConn, timeout uint) (replies map[*net.UDPAddr][][]byte, err error) {

	conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeout)))
	replies = make(map[*net.UDPAddr][][]byte, 0)
	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if err.(net.Error).Timeout() {
				return replies, nil
			}
		}
		if replies[src] == nil {
			replies[src] = make([][]byte, 0)
		}
		replies[src] = append(replies[src], buffer[0:numBytes])
	}
}

func read(conn *net.UDPConn) (buffer []byte, src *net.UDPAddr, err error) {

	conn.SetReadBuffer(maxDatagramSize)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	buffer = make([]byte, maxDatagramSize)
	numBytes, src, err := conn.ReadFromUDP(buffer)
	buffer = buffer[0:numBytes]

	return
}
