package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/waynz0r/gree-ac-client/pkg/client"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}
}

func main() {

	responses, err := client.Scan("192.168.255.255:7000", 1000)
	checkError(err)

	for _, reply := range responses {
		client, err := client.NewClient(reply.Src.String(), reply.Response.MAC)
		checkError(err)
		status, err := client.Status()
		checkError(err)
		spew.Dump(status)
		break
	}

	// for _, r := range responses {
	// 	fmt.Printf("Response from %s: %s/%s\n", r.Src, r.Response.Name, r.Response.MAC)
	// 	// spew.Dump(r.UnPack())
	// }

	// CheckError(err)
	// resp, src, err := client.Scan("192.168.255.255:7000")

	// fmt.Printf("Response from %s: %s/%s\n", src, resp.Name, resp.MAC)

	// key, err := client.Bind(src.String(), resp.MAC)
	// CheckError(err)

	// fmt.Printf("Bind OK, specific key: %s\n", key)

	// client := client.NewClient(src.String(), resp.MAC, key)

	// client.PowerOn()
	// client.SetTemperature(18)
	// client.SetTemperature(20)
	// client.SwingMode(2)
}
