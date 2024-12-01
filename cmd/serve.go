package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	tplt "go-image-annotator/templates"
	"log"
	g "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port     int
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			serve(port)
		},
	}
)

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 80, "port to serve on")
}

func serve(port int) {
	mux := http.NewServeMux()

	items := []string{"Super", "Duper", "Nice"}
	mux.Handle("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		return tplt.HomePage(items), nil
	}))

	server := &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: mux}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	log.Printf("Starting server on port %v...\n", port)
	log.Printf("API docs URL: <root>:%v/docs\n", port)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("server stopping...")
	defer cancel()

	log.Fatal(server.Shutdown(ctx))

}
