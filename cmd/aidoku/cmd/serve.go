package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/beerpiss/aidoku-cli/internal/build"
	"github.com/beerpiss/aidoku-cli/internal/logcat"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] \"%s %s\" \"%s\"", time.Now().UTC().Format(time.RFC3339), r.Method, r.URL, r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}

var serveCmd = &cobra.Command{
	Use:           "serve <FILES>",
	Short:         "Build a source list and serve it on the local network",
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

		output, _ := cmd.Flags().GetString("output")
		port, _ := cmd.Flags().GetString("port")

		os.RemoveAll(output)
		var fileList []string
		for _, arg := range args {
			files, err := filepath.Glob(arg)
			if err != nil {
				color.Red("error: invalid glob pattern %s", arg)
				continue
			}
			fileList = append(fileList, files...)
		}
		if len(fileList) == 0 {
			return errors.New("no files given")
		}
		os.MkdirAll(output, os.FileMode(0644))
		os.MkdirAll(output+"/icons", os.FileMode(0644))
		os.MkdirAll(output+"/sources", os.FileMode(0644))
		build.BuildSource(fileList, output)

		http.Handle("/", http.FileServer(http.Dir(output)))
		fmt.Println("Listening on these addresses:")
		logcat.PrintAddresses(port)

		http.ListenAndServe(":"+port, logRequest(http.DefaultServeMux))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8080", "The port to broadcast the source list on")
	serveCmd.Flags().StringP("output", "o", "public", "The source list folder")

	serveCmd.MarkZshCompPositionalArgumentFile(1, "*.aix")
	serveCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"aix"}, cobra.ShellCompDirectiveFilterFileExt
	}
}
