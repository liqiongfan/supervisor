package supervisor

import (
	`net`
	`net/http`
	`google.golang.org/grpc`
	`errors`
	`flag`
	`os`
)


type GRPCServer struct {
	Server
	Entry      func([]*grpc.Server, []net.Listener)
	
	// internal properties
	server   []*grpc.Server
}



func (h *GRPCServer) restart(w http.ResponseWriter, r *http.Request) {
	h.commonRestart(w, r)
	
	go func(h *GRPCServer) {
		for i, s := range h.server {
			s.GracefulStop()
			h.cs[i] <- 1
		}
	}(h)
}


func (h *GRPCServer) notify() (err error) {
	
	if h.Config.Addr == "" {
		return errors.New(`ListenConfig's addr empty`)
	}
	
	go func( h *GRPCServer) {
		
		mux := http.NewServeMux()
		mux.HandleFunc(`/-/reload`, h.restart)
		
		s := http.Server{Addr: h.Config.Addr, Handler: mux}
		err = s.ListenAndServe()
		
	}(h)
	
	return nil
}


func (h *GRPCServer) Run() (err error) {

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
		
		h.server = append(h.server, &grpc.Server{})
		
		h.cs = append(h.cs, make(chan int))
		
	}
	
	err = h.notify()
	h.Entry(h.server, h.listener)
	
	for _, c := range h.cs {
		<- c
	}
	
	return nil
	
}