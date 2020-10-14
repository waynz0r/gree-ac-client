package client

import (
	"encoding/base64"
	"encoding/json"
)

type RequestInterface interface {
	GetMAC() string
	GetType() string
	GetEncrypted([]byte) ([]byte, error)
}

type Reply struct {
	ClientID string `json:"cid"`
	I        uint   `json:"i"`
	Pack     string `json:"pack"`
	Type     string `json:"t"`
	MAC      string `json:"tcid,omitempty"`
	UID      uint   `json:"uid"`
}

type Message struct {
	ClientID string `json:"cid"`
	I        uint   `json:"i"`
	Pack     string `json:"pack"`
	Type     string `json:"t"`
	MAC      string `json:"tcid,omitempty"`
	UID      uint   `json:"uid"`
}

type BaseRequest struct {
	MAC  string `json:"mac"`
	Type string `json:"t"`
}

type BaseResponse struct {
	MAC  string `json:"mac"`
	Type string `json:"t"`
}

func (br BaseRequest) GetMAC() string {

	return br.MAC
}

func (br BaseRequest) GetType() string {

	return br.Type
}

func (br BaseRequest) GetEncrypted(key []byte) (encrypted []byte, err error) {

	r, err := json.Marshal(br)
	if err != nil {
		return
	}

	encrypted, err = encrypt(r, key)

	return
}

func (r *Reply) DecryptPack(key string) (decrypted string, err error) {

	if r.Type != "pack" {
		return
	}

	cipherText, err := base64.StdEncoding.DecodeString(r.Pack)
	if err != nil {
		return
	}

	plaintext, err := decrypt(cipherText, []byte(key))
	if err != nil {
		return
	}

	decrypted = string(plaintext)

	return
}
