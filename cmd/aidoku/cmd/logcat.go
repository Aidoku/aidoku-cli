package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Aidoku/aidoku-cli/internal/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func Logcat(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Fprintf(w, "Method not supported\n")
	} else {
		buf := new(strings.Builder)
		io.Copy(buf, req.Body)

		items := strings.Split(buf.String(), "] [")
		t, err := time.Parse("01/02 03:04:05.999", strings.TrimLeft(items[0], "["))
		if err == nil {
			t = time.Date(time.Now().Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
			items[0] = fmt.Sprintf("[%s", t.Format("2006-01-02T15:04:05.999999999"))
		}
		log := strings.Join(items, "] [")

		if strings.Contains(log, "[ERROR]") {
			color.Red(log)
		} else if strings.Contains(log, "[WARN]") {
			color.Yellow(log)
		} else if strings.Contains(log, "[DEBUG]") {
			color.HiBlack(log)
		} else {
			fmt.Println(log)
		}
	}
}

var logcatCmd = &cobra.Command{
	Use:           "logcat",
	Short:         "Log streaming",
	Version:       rootCmd.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			os.Exit(0)
		}()

		if ForceColor {
			color.NoColor = false
		}

		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetString("port")
		http.HandleFunc("/", Logcat)

		fmt.Println("Listening on these addresses:")
		if address == "0.0.0.0" {
			common.PrintAddresses(port)
		} else {
			color.Green("    http://%s:%s", address, port)
		}

		http.ListenAndServe(address+":"+port, nil)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logcatCmd)
	logcatCmd.Flags().StringP("address", "a", "0.0.0.0", "Address to listen to logs on")
	logcatCmd.Flags().StringP("port", "p", "9000", "The port to listen to logs on")
}
