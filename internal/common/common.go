package common

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fatih/color"
	"github.com/spf13/cast"
)

func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
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
	color.Green("    http://127.0.0.1:%s\n", port)
}

func GeneratePng(location string) error {
	img, err := os.Create(location)
	if err != nil {
		color.Red("error: Couldn't write icon file %s: %s", location, err.Error())
		return err
	}
	transparent, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg==")
	io.Copy(img, bytes.NewReader(transparent))
	img.Sync()
	img.Close()
	return nil
}

func ProcessGlobs(globs []string) []string {
	var fileList []string
	for _, arg := range globs {
		base, pattern := doublestar.SplitPattern(filepath.ToSlash(arg))
		fsys := os.DirFS(base)
		files, err := doublestar.Glob(fsys, pattern)
		if err != nil {
			color.Red("error: invalid glob pattern %s", arg)
			continue
		}
		for _, file := range files {
			fileList = append(fileList, base+"/"+file)
		}
	}
	return fileList
}

func ToDurationE(v any) (time.Duration, error) {
	if n := cast.ToInt(v); n > 0 {
		return time.Duration(n) * time.Millisecond, nil
	}
	d, err := time.ParseDuration(cast.ToString(v))
	if err != nil {
		return 0, fmt.Errorf("cannot convert %v to time.Duration", v)
	}
	return d, nil
}
