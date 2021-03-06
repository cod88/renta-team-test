package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	_ "github.com/mattn/go-sqlite3"
)

type StoreConfig struct {
	DbFilename string
}

type WholeConfig struct {
	StoreConfig StoreConfig
}

type NewsRecord struct {
	Id    string
	Title string
	Date  string
}

var Config StoreConfig
var DB *sql.DB

func init() {
	fmt.Println("Init Store server")
	configure()

	db, err := sql.Open("sqlite3", "../"+Config.DbFilename)

	if err != nil {
		fmt.Printf("Database: %+v\n", err)
	} else {
		DB = db
	}
}

func configure() {
	var wCfg WholeConfig

	execFile, _ := os.Executable()
	approot := filepath.Dir(filepath.Dir(execFile))

	if _, err := toml.DecodeFile(approot+"/config/config.toml", &wCfg); err != nil {
		fmt.Println("We have an error on get StoreConfig config. ", err)
	}
	Config = wCfg.StoreConfig
	wCfg = WholeConfig{}
}

func GetNews(id string) string {

	var rec NewsRecord

	rows, err := DB.Query("SELECT id, title, date FROM news WHERE id=?", id)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "error"
	}

	defer rows.Close()
	rows.Next()
	rows.Scan(&rec.Id, &rec.Title, &rec.Date)

	var data string
	if rec.Id == "" {
		data = "News does not exists"
	} else {
		unmarsheled, _ := json.Marshal(rec)
		data = string(unmarsheled)
	}

	fmt.Printf("%+v\n", data)
	return data
}
