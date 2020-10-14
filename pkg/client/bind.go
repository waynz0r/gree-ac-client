package client

import (
	"encoding/json"
	"fmt"
)

type BindRequest struct {
	BaseRequest

	UID uint `json:"uid"`
}

type BindResponse struct {
	BaseRequest

	Key string `json:"key"`
}

type BindCommand struct {
}

type BindReply struct {
	Raw      []byte
	Reply    *Reply
	Response *BindResponse
}

func (br BindRequest) GetEncrypted(key []byte) (encrypted []byte, err error) {

	r, err := json.Marshal(br)
	if err != nil {
		return
	}

	encrypted, err = encrypt(r, key)

	return
}

func (c *Client) Bind() (key string, err error) {

	req, err := newMessage(BindRequest{
		BaseRequest: BaseRequest{Type: "bind", MAC: c.MAC},
	}, defaultKey, false)

	resp, err := send(c.Address, req)

	message, err := parseReply(resp, defaultKey)
	if err != nil {
		return
	}

	var m BindResponse
	err = json.Unmarshal([]byte(message.Pack), &m)
	if err != nil {
		return
	}

	if m.Type != "bindok" {
		err = fmt.Errorf("Invalid bind response from %s", c.Address)
		return
	}

	key = m.Key

	return
}
