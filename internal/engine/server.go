package engine

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/golang/glog"
	"github.com/livepeer/swarm-chaos/internal/model"
)

const bindAddress = "0.0.0.0:7933"

type (
	// Server serves API
	Server struct {
		scheduler *Scheduler
	}

	scheduleTaskRequest struct {
		IntMin      string `json:"int_min,omitempty"`
		IntMax      string `json:"int_max,omitempty"`
		FilterKey   string `json:"filter_key,omitempty"`
		FilterValue string `json:"filter_value,omitempty"`
	}

	versionResponse struct {
		Version  string `json:"version,omitempty"`
		OS       string `json:"os,omitempty"`
		ARCH     string `json:"arch,omitempty"`
		Compiler string `json:"compiler,omitempty"`
		Runtime  string `json:"runtime,omitempty"`
	}
)

// NewServer returns a new Server
func NewServer(scheduler *Scheduler) *Server {
	return &Server{scheduler: scheduler}
}

// StartServer start serving
// blocks until server is stopped
func (srv *Server) StartServer() error {
	mux := srv.webServerHandlers(bindAddress)
	s := &http.Server{
		Addr:    bindAddress,
		Handler: mux,
	}

	glog.Info("Web server listening on ", bindAddress)
	return s.ListenAndServe()
}

func (srv *Server) webServerHandlers(bindAddr string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/schedule_task", func(w http.ResponseWriter, r *http.Request) {
		srv.handleScheduleTask(w, r)
	})
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		// ss.handleStats(w, r)
	})
	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		srv.handleStop(w, r)
	})
	mux.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		srv.handleClear(w, r)
	})
	mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		srv.handleStart(w, r)
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		srv.handleVersion(w, r)
	})
	return mux
}

// Schedule task
func (srv *Server) handleScheduleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	str := &scheduleTaskRequest{}
	err = json.Unmarshal(b, str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	glog.Infof("Got schedule task request %+v.", *str)

	err = srv.scheduler.ScheduleTask(str.IntMin, str.IntMax, model.OperationTypeDestroy, str.FilterKey, str.FilterValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Start scheduled tasks
func (srv *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	glog.Info("Got start request.")
	err := srv.scheduler.StartTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Stop currently running tasks
func (srv *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	glog.Info("Got stop request.")
	srv.scheduler.StopTasks()
	w.WriteHeader(http.StatusOK)
}

// Clear scheduled tasks
func (srv *Server) handleClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	glog.Info("Got clear tasks request.")
	srv.scheduler.ClearTasks()
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	resp := versionResponse{
		Version:  model.SwarmChaosVersion,
		Compiler: runtime.Compiler,
		Runtime:  runtime.Version(),
		ARCH:     runtime.GOARCH,
		OS:       runtime.GOOS,
	}
	respB, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(respB)
}
