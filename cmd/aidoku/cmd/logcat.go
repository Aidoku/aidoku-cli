package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/beerpiss/aidoku-cli/internal/logcat"
	"github.com/spf13/cobra"
)

var logcatCmd = &cobra.Command{
	Use:           "logcat",
	Short:         "Log streaming",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			os.Exit(0)
		}()

		port, _ := cmd.Flags().GetString("port")
		http.HandleFunc("/", logcat.Logcat)

		fmt.Println("Listening on these addresses:")
		logcat.PrintAddresses(port)

		http.ListenAndServe(":"+port, nil)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logcatCmd)
	logcatCmd.Flags().StringP("port", "p", "9000", "The port to listen to logs on")
}
