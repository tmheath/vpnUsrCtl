package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	router "github.com/swoga/go-routeros"
	"net"
	"os"
)

func main() {
	router_addr := ""
	router_user := ""
	router_password := ""
	a := app.New()
	w := a.NewWindow("VPN")
	conn, err := net.Dial("udp", router_addr)
	if err != nil {
		w.SetContent(widget.NewLabel("Error, Network down"))
		w.ShowAndRun()
		os.Exit(4)
	}
	local_addr := conn.LocalAddr().(*net.UDPAddr).IP
	drop_cmd := fmt.Sprintf("/ip/firewall/mangle add chain=forward src-address=%s action=mark-packet new-packet-mark=NOVPN", local_addr)
	revert_cmd := fmt.Sprintf("")
	client, err := router.Dial(router_addr, router_user, router_password)
	if err != nil {
		w.SetContent(widget.NewLabel("Error, failed to connect to router."))
		w.ShowAndRun()
		os.Exit(1)
	}
	_, err = client.Run(drop_cmd)
	if err != nil {
		w.SetContent(widget.NewLabel("Error, failed to set vpn status"))
		w.ShowAndRun()
		os.Exit(2)
	}
	w.SetContent(widget.NewLabel("Success, local traffic off vpn"))
	w.ShowAndRun()
	_, err = client.Run(revert_cmd)
	if err != nil {
		ew := a.NewWindow("VPN - Error")
		ew.SetContent(widget.NewLabel("Error, VPN Remains Inactive"))
		ew.ShowAndRun()
		os.Exit(3)
	}
}

