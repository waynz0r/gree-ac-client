package client

import (
	"encoding/json"
	"fmt"
	"net"
)

// ScanRequest represents a scan request
type ScanRequest BaseRequest

// ScanResponse represents a scan response
type ScanResponse struct {
	BaseResponse
	ClientID string `json:"cid"`
	Model    string `json:"model"`
	Name     string `json:"name"`
	Version  string `json:"ver"`
	Lock     uint   `json:"lock"`
}

// ScanCommand represents a Scan command
type ScanCommand struct {
	broadcast string
	timeout   uint
}

// ScanReply holds data replied to a ScanCommand
type ScanReply struct {
	Src      *net.UDPAddr
	Raw      []byte
	Reply    *Reply
	Response *ScanResponse
}

// Scan scans the given broadcast address for compatible equipments
func Scan(broadcast string, timeout uint) (responses []*ScanReply, err error) {

	cmd := &ScanCommand{
		broadcast: broadcast,
		timeout:   timeout,
	}

	return cmd.scan()
}

func (cmd *ScanCommand) scan() (responses []*ScanReply, err error) {

	conn, err := listen(cmd.broadcast)
	if err != nil {
		return
	}

	addr, err := net.ResolveUDPAddr("udp", cmd.broadcast)
	if err != nil {
		return
	}

	req, err := json.Marshal(ScanRequest{Type: "scan"})
	if err != nil {
		return
	}

	conn.WriteTo(req, addr)

	_replies, err := readForAWhile(conn, cmd.timeout)
	if err != nil {
		return
	}

	return cmd.parseReplies(_replies)
}

func (cmd *ScanCommand) parseReplies(replies map[*net.UDPAddr][][]byte) (parsedReplies []*ScanReply, err error) {

	parsedReplies = make([]*ScanReply, 0)
	for src, rs := range replies {
		for _, r := range rs {
			reply, err := cmd.parseReply(src, r)
			if err != nil {
				return parsedReplies, err
			}
			parsedReplies = append(parsedReplies, reply)
		}
	}

	return
}

func (cmd *ScanCommand) parseReply(src *net.UDPAddr, reply []byte) (parsedReply *ScanReply, err error) {

	parsedReply = &ScanReply{
		Src: src,
		Raw: reply,
	}

	r, err := parseReply(reply, defaultKey)
	if err != nil {
		return
	}
	parsedReply.Reply = r

	pack, err := parsedReply.unpack()
	if err != nil {
		return
	}
	parsedReply.Response = pack

	return
}

func (sr *ScanReply) unpack() (r *ScanResponse, err error) {

	err = json.Unmarshal([]byte(sr.Reply.Pack), &r)

	if r.Type != "dev" {
		err = fmt.Errorf("Invalid scan response from %s", sr.Src)
		return
	}

	return
}
