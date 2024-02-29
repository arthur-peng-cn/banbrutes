package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arthur/banbrutes/config"
	"github.com/zngw/log"
)

// Server service
type Service struct {
	// Dispatch connections to different handlers listen on same port
	cfg config.Conf
}

func NewService(cfg config.Conf) (svr *Service, err error) {

	svr = &Service{
		cfg: cfg,
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Trace("link", "GET访问")
		data := []map[string]string{{"提示": "非法访问"}}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(data)
		return
	}

	// log.Println("POST访问")
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Trace("link", "%v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	isAllowSSH := handleMsg(data)
	var responseData map[string]interface{}
	if isAllowSSH {
		responseData = map[string]interface{}{"reject": false, "unchange": true}
	} else {
		responseData = map[string]interface{}{"reject": true, "reject_reason": "invalid user"}
	}

	log.Trace("link", "%s", responseData)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func handler2(w http.ResponseWriter, r *http.Request) {

	responseData := map[string]interface{}{"reject": false, "unchange": true}
	log.Trace("link", "%s", responseData)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func (srv *Service) Run() {
	//
	for _, filter := range srv.cfg.Filters {
		// for _, reg := range filter.RegFilters {
		logMonitorServer(filter.LogFile, filter.RegFilters)

	}
	if strings.TrimSpace(srv.cfg.ListenAddr) != "" {
		http.HandleFunc("/handler", handler)
		http.HandleFunc("/info", handler2)
		log.Info("sys", "%v", http.ListenAndServe(srv.cfg.ListenAddr, nil))
	}
}
