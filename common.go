package supervisor

import (
	`flag`
	`net`
	`os`
	`os/exec`
	`errors`
	`net/http`
	`fmt`
)

var (
	rerun = flag.Bool(`notify`, false, `Program should run again and exit the current program`)
)


type ListenConfig struct {
	Addr string
}


type Server struct {
	ListenAddr []string
	Config     ListenConfig
	
	// internal properties
	listener   []net.Listener
	cs         []chan int
}


func (h *Server) reload() (err error) {
	var argExist = false
	for _, arg := range os.Args {
		if arg == `-notify` {
			argExist = true
			break
		}
	}
	if !argExist {
		os.Args = append(os.Args, `-notify`)
	}
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	
	var fs []*os.File
	for _, listener := range h.listener {
		l, ok := listener.(*net.TCPListener)
		if !ok {
			return errors.New(`listener not TCPListener`)
		}
		f, err := l.File()
		if err != nil {
			return err
		}
		fs = append(fs, f)
	}
	cmd.ExtraFiles = fs
	err = cmd.Start()
	
	return nil
}


func (h *Server) commonRestart(w http.ResponseWriter, r *http.Request) {
	var err error
	_, _ = w.Write([]byte("success\n"))
	
	err = h.reload()
	if err != nil {
		_ = fmt.Errorf("gracefully shutdown error: %v\n", err)
		return
	}
}