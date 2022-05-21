package logcat

import (
	"fmt"
	"io"
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
