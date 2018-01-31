package debug

import (
	"net/http/pprof"

	"github.com/gorilla/mux"
)

func PProfHandle(r *mux.Router) {
	r.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
	r.PathPrefix("/debug/pprof/cmdline").HandlerFunc(pprof.Cmdline)
	r.PathPrefix("/debug/pprof/profile").HandlerFunc(pprof.Profile)
	r.PathPrefix("/debug/pprof/symbol").HandlerFunc(pprof.Symbol)
	r.PathPrefix("/debug/pprof/trace").HandlerFunc(pprof.Trace)
}
