package client

const (
	defaultKey = "a3K8Bx%2r8Y7#xDh"
)

type Client struct {
	Address string
	MAC     string
	Key     string
}

func NewClient(address string, MAC string) (client *Client, err error) {

	client = &Client{
		Address: address,
		MAC:     MAC,
	}

	key, err := client.Bind()
	if err != nil {
		return
	}

	client.Key = key

	return
}
