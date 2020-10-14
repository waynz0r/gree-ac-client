package client

import (
	"encoding/json"
	"fmt"
)

type DeviceStatus map[string]int

type StatusRequest struct {
	BaseRequest

	Cols []string `json:"cols"`
}

type StatusResponse struct {
	BaseRequest

	Cols []string `json:"cols"`
	Data []int    `json:"dat"`
}

func (c *Client) Status() (status DeviceStatus, err error) {

	req, err := newMessage(StatusRequest{
		BaseRequest: BaseRequest{Type: "status", MAC: c.MAC},
		Cols:        []string{"Pow", "Mod", "SetTem", "WdSpd", "Air", "Blo", "Health", "SwhSlp", "Lig", "SwingLfRig", "SwUpDn", "Quiet", "Tur", "StHt", "TemUn", "HeatCoolType", "TemRec", "SvSt"},
	}, c.Key, true)

	resp, err := send(c.Address, req)

	message, err := parseReply(resp, c.Key)
	if err != nil {
		return
	}

	var s StatusResponse
	err = json.Unmarshal([]byte(message.Pack), &s)
	if err != nil {
		return
	}

	if s.Type != "dat" {
		err = fmt.Errorf("Invalid status response from %s", c.Address)
		return
	}

	status = make(DeviceStatus, 0)
	for i, name := range s.Cols {
		status[name] = s.Data[i]
	}

	return
}

func (sr StatusRequest) GetEncrypted(key []byte) (encrypted []byte, err error) {

	r, err := json.Marshal(sr)
	if err != nil {
		return
	}

	encrypted, err = encrypt(r, key)

	return
}
