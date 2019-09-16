//+build !windows

package grace

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

)

type Grace struct {
	Server   *http.Server
	Listener net.Listener
}

func New(addr string, handler http.Handler, graceful bool) (*Grace, error) {
	var ln net.Listener
	var err error
	if graceful {
		log.Println("Listening to existing file descriptor 3.")
		f := os.NewFile(3, "")
		ln, err = net.FileListener(f)
	} else {
		log.Println("Listening on a new file descriptor.")
		ln, err = net.Listen("tcp", addr)
	}
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Handler:     handler,
	}
	return &Grace{
		Server:   s,
		Listener: ln,
	}, nil
}

func (g *Grace) Start() {
	go func() {
		if err := g.Server.Serve(g.Listener); err != nil {
			log.Println(err)
		} else {
			log.Println("Shutdown old server...")
		}
	}()

	g.watchSign()
}

func (g *Grace) reload() error {
	tl, ok := g.Listener.(*net.TCPListener)
	if !ok {
		return fmt.Errorf("%s", "listener is not tcp listener")
	}
	f, err := tl.File()
	if err != nil {
		return err
	}

	args := []string{"-graceful"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// put socket FD at the first entry
	cmd.ExtraFiles = []*os.File{f}
	return cmd.Start()
}

func (g *Grace) watchSign() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	for {
		sig := <-ch
		log.Println("watchSign", sig.String())
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// stop
			log.Println("watchSign", "stop server...")
			signal.Stop(ch)
			if err := g.Server.Shutdown(); err != nil {
				log.Println("watchSign", err)
			}
			return
		case syscall.SIGUSR2:
			// reload
			log.Println("watchSign", "reload server...")
			err := g.reload()
			if err != nil {
				log.Println("watchSign", err)
			}
			if err := g.Server.Shutdown(); err != nil {
				log.Println("watchSign", err)
			}
			return
		}
	}
}
