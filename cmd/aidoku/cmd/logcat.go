package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Aidoku/aidoku-cli/internal/common"
	"github.com/Aidoku/aidoku-cli/internal/logcat"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

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
		http.HandleFunc("/", logcat.Logcat)

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
