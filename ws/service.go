package ws

import (
	"FamilyWatch/conf"
	"flag"
	"log"
	"net/http"
)

func Start() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		serveWss(hub, w, r)
	})
	err := http.ListenAndServe(conf.Conf.WsAddr, nil)
	//err := http.ListenAndServeTLS(conf.Conf.WsAddr,
	//	conf.Conf.CertFile,
	//	conf.Conf.KeyFile, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (req Request) Process() (resp Respond) {
	switch req.Op {
	case 1:
	case 2:
	case 3:
	default:

	}
	return resp
}
