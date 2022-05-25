package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aidoku/aidoku-cli/internal/build"
	"github.com/Aidoku/aidoku-cli/internal/common"
	"github.com/fatih/color"
	"github.com/felixge/httpsnoop"
	"github.com/spf13/cobra"
)

var (
	cyan = color.New(color.FgCyan).SprintFunc()
	red  = color.New(color.FgRed).SprintFunc()
)

var serveCmd = &cobra.Command{
	Use:           "serve <FILES>",
	Short:         "Build a source list and serve it on the local network",
	Version:       rootCmd.Version,
	Args:          cobra.MinimumNArgs(1),
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
		output, _ := cmd.Flags().GetString("output")
		port, _ := cmd.Flags().GetString("port")

		build.BuildWrapper(args, output)

		fmt.Println("Listening on these addresses:")
		if address == "0.0.0.0" {
			common.PrintAddresses(port)
		} else {
			color.Green("    http://%s:%s", address, port)
		}
		fmt.Println("Hit CTRL-C to stop the server")

		handler := http.FileServer(http.Dir(output))
		http.Handle("/", handler)
		wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timestamp := time.Now().UTC().Format(time.RFC3339)
			method := r.Method
			url := r.URL
			userAgent := r.UserAgent()
			fmt.Printf("[%s] \"%s %s\" \"%s\"\n", timestamp, cyan(method), cyan(url), userAgent)

			m := httpsnoop.CaptureMetrics(handler, w, r)
			timestamp = time.Now().UTC().Format(time.RFC3339)
			statusCode := m.Code
			if statusCode != http.StatusOK {
				fmt.Printf("[%s] \"%s %s\" Error (%s): \"%s\"\n", timestamp, red(method), red(url), red(statusCode), red(http.StatusText(statusCode)))
			}
		})
		return http.ListenAndServe(address+":"+port, wrappedHandler)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("address", "a", "0.0.0.0", "Address to broadcast source list")
	serveCmd.Flags().StringP("port", "p", "8080", "The port to broadcast the source list on")
	serveCmd.Flags().StringP("output", "o", "public", "The source list folder")

	serveCmd.MarkZshCompPositionalArgumentFile(1, "*.aix")
	serveCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"aix"}, cobra.ShellCompDirectiveFilterFileExt
	}
}
