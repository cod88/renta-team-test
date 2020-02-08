package requestserver

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	mux "github.com/gorilla/mux"

	"github.com/BurntSushi/toml"
)

type RequestServerConfig struct {
	Port string
	Host string
}

type WholeConfig struct {
	RequestServerConfig RequestServerConfig
}

var Config RequestServerConfig

func init() {
	fmt.Println("Init request server")
	configure()
}

func configure() {
	var wCfg WholeConfig

	execFile, _ := os.Executable()
	approot := filepath.Dir(filepath.Dir(execFile))

	if _, err := toml.DecodeFile(approot+"/config/config.toml", &wCfg); err != nil {
		fmt.Println("We have an error on get RequestServerConfig config. ", err)
	}
	Config = wCfg.RequestServerConfig
	wCfg = WholeConfig{}
}

func RunServer() {
	router := mux.NewRouter()
	router.HandleFunc("/news/{id}", getNews).Methods("GET")

	http.ListenAndServe(Config.Host+":"+Config.Port, router)
}

func getNews(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "That news: %s !", vars["id"])
}
