# supervisor
Golang http&amp;grpc server for gracefully shutdown like nginx -s reload
if you want a server which would be restart without stop service, you shall choise supervisor

# demo

HTTP server

```golang

package main

import "github.com/liqiongfan/supervisor"


func main() {

    h := &supervisor.HTTPServer{
        Server: supervisor.Server{
            ListenAddr: []string{`:9091`},
            Config: supervisor.ListenConfig{ Addr: ":8088" },
        },
        Entry: Main,
    }

    err := h.Run()
    if err != nil {
        panic(err)
    }
}


func Main(srv []*http.Server, l []net.Listener) {

    http.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request){
        _, _ = w.Write([]byte("OK\n"))
    })

    err := srv[0].Serve(l[0])
    if err != nil {
        panic(err)
    }

}



```