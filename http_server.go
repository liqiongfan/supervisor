package supervisor

import (
	`net`
	`net/http`
	`errors`
	`flag`
	`os`
	`context`
	`fmt`
)


type HTTPServer struct {
	Server
	Entry      func([]*http.Server, []net.Listener)
	
	// internal properties
	server     []*http.Server
}


func (h *HTTPServer) restart(w http.ResponseWriter, r *http.Request) {
	var err error
	
	h.commonRestart(w, r)
	
	go func(h *HTTPServer) {
		for i, s := range h.server {
			err = s.Shutdown(context.Background())
			if err != nil {
				_ = fmt.Errorf("gracefully shutdown error: %v\n", err)
				continue
			}
			h.cs[i] <- 1
		}
	}(h)
}


func (h *HTTPServer) notify() (err error) {
	
	if h.Config.Addr == "" {
		return errors.New(`ListenConfig's addr empty`)
	}
	
	go func( h *HTTPServer) {
		
		mux := http.NewServeMux()
		mux.HandleFunc(`/-/reload`, h.restart)
		
		s := http.Server{Addr: h.Config.Addr, Handler: mux}
		err = s.ListenAndServe()
		
	}(h)
	
	return nil
}


func (h *HTTPServer) Run() (err error) {

	var listener net.Listener
	
	flag.Parse()
	
	for i, addr := range h.ListenAddr {
		
		if *rerun {
			f := os.NewFile(uintptr(3 + i), addr)
			listener, err = net.FileListener(f)
		} else {
			listener, err = net.Listen(`tcp`, addr)
		}
		
		if err != nil {
			return
		}
		h.listener = append(h.listener, listener)
		
		h.server = append(h.server, &http.Server{})
		
		h.cs = append(h.cs, make(chan int))
		
	}
	
	err = h.notify()
	h.Entry(h.server, h.listener)
	
	for _, c := range h.cs {
		<- c
	}
	
	return nil
	
}