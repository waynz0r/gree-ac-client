package client

import "encoding/json"

type ControlRequest struct {
	BaseRequest

	Options []string `json:"opt"`
	Values  []uint   `json:"p"`
}

func (cr ControlRequest) GetEncrypted(key []byte) (encrypted []byte, err error) {

	r, err := json.Marshal(cr)
	if err != nil {
		return
	}

	encrypted, err = encrypt(r, key)

	return
}

func (c *Client) SetMode(mode uint) (bool, error) {

	return c.cmd([]string{"Mod"}, []uint{mode})
}

func (c *Client) PowerOff() (bool, error) {

	return c.cmd([]string{"Pow"}, []uint{0})
}

func (c *Client) PowerOn() (bool, error) {

	return c.cmd([]string{"Pow"}, []uint{1})
}

func (c *Client) SetTemperature(temperature uint) (bool, error) {

	return c.cmd([]string{"TemUn", "SetTem"}, []uint{0, temperature})
}

func (c *Client) LightOff() (bool, error) {

	return c.cmd([]string{"Lig"}, []uint{0})
}

func (c *Client) LightOn() (bool, error) {

	return c.cmd([]string{"Lig"}, []uint{1})
}

func (c *Client) SwingMode(mode uint) (bool, error) {

	return c.cmd([]string{"SwUpDn"}, []uint{mode})
}

func (c *Client) cmd(options []string, values []uint) (bool, error) {

	req, err := newMessage(ControlRequest{
		BaseRequest: BaseRequest{Type: "cmd", MAC: c.MAC},
		Options:     options,
		Values:      values,
	}, c.Key, true)
	if err != nil {
		return false, err
	}

	send(c.Address, req)

	return true, nil
}
