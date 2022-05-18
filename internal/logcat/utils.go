package logcat

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

func Logcat(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Fprintf(w, "Method not supported\n")
	} else {
		buf := new(strings.Builder)
		io.Copy(buf, req.Body)
		log := buf.String()
		if strings.Contains(log, "error") {
			color.Red(log)
		} else if strings.Contains(log, "warning") {
			color.Yellow(log)
		} else {
			fmt.Println(log)
		}
	}
}

func PrintAddresses(port string) {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				if isPrivateIP(ip) {
					color.Green("    http://%s:%s\n", ip.String(), port)
				} else if !strings.Contains(ip.String(), ":") {
					fmt.Printf("    http://%s:%s\n", ip.String(), port)
				}
			}
		}
	}
}

func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
