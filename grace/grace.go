//+build !windows

package grace

import (
	"log"
	"net"
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
		log.Warn("grace", zap.String("new", "Listening to existing file descriptor 3."))
		f := os.NewFile(3, "")
		ln, err = net.FileListener(f)
	} else {
		log.Warn("grace", zap.String("new", "Listening on a new file descriptor."))
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
			log.Error("start", zap.String("server.Serve err", err.Error()))
		} else {
			log.Warn("start", zap.String("server.Serve", "Shutdown old server..."))
		}
	}()

	g.watchSign()
}

func (g *Grace) reload() error {
	tl, ok := g.Listener.(*net.TCPListener)
	if !ok {
		return errors.New("listener is not tcp listener")
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
		log.Info("watchSign", zap.String("signal", sig.String()))
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// stop
			log.Info("watchSign", zap.String("stop", "stop server..."))
			signal.Stop(ch)
			if err := g.Server.Shutdown(); err != nil {
				log.Error("watchSign", zap.String("server shutdown err", err.Error()))
			}
			return
		case syscall.SIGUSR2:
			// reload
			log.Info("watchSign", zap.String("reload", "reload server..."))
			err := g.reload()
			if err != nil {
				log.Error("watchSign", zap.String("reload server err", err.Error()))
			}
			if err := g.Server.Shutdown(); err != nil {
				log.Error("watchSign", zap.String("server shutdown err", err.Error()))
			}
			return
		}
	}
}
