package messenger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"store"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessengerConfig struct {
	Username  string
	NewsTopic string
	BrokerUrl string
}

type WholeConfig struct {
	MessengerConfig MessengerConfig
}

type TQueryNews struct {
	Id string `json:"newsid"`
}

var Config MessengerConfig
var mqttClient mqtt.Client

func init() {
	fmt.Println("Init messenger")
	configure()
}

func configure() {
	var wCfg WholeConfig

	execFile, _ := os.Executable()
	approot := filepath.Dir(filepath.Dir(execFile))

	if _, err := toml.DecodeFile(approot+"/config/config.toml", &wCfg); err != nil {
		fmt.Println("We have an error on get MessengerConfig config. ", err)
	}
	Config = wCfg.MessengerConfig
	wCfg = WholeConfig{}
}

func QueryNews(id string) error {
	var qn TQueryNews
	qn.Id = id
	jsoned, err := json.Marshal(qn)

	if err != nil {
		return err
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	if token := client.Publish(Config.NewsTopic, 0, false, string(jsoned)); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Println("Query for news:" + string(jsoned) + " is sent")
	return nil
}

func WaitAnswerForNews(id string) (data string, err error) {
	client, err := getClient()
	if err != nil {
		return "", err
	}
	var wg sync.WaitGroup
	nid := id
	var answer string
	wg.Add(1)
	client.Subscribe(Config.NewsTopic+"/"+id, 0, func(client mqtt.Client, message mqtt.Message) {
		fmt.Println("News " + nid + " waited" + string(message.Payload()))
		answer = string(message.Payload())
		wg.Done()
	})

	fmt.Println("Wait for " + id)

	wg.Wait()
	fmt.Println("Awaited for " + id + ":" + answer)
	return answer, nil
}

func getClient() (mqtt.Client, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("tcp://%s", Config.BrokerUrl))
	options.SetUsername(Config.Username)

	mqttClient := mqtt.NewClient(options)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	} else {
		return mqttClient, nil
	}
}

func WaitQueryForNews() error {
	client, err := getClient()

	if err != nil {
		return err
	}

	client.Subscribe(Config.NewsTopic, 0, func(client mqtt.Client, message mqtt.Message) {
		fmt.Println("Query for news " + string(message.Payload()))
		var qn TQueryNews
		var result string

		err := json.Unmarshal(message.Payload(), &qn)
		if err != nil {
			result = "error"
			return
		}

		fmt.Printf("%+v\n", qn)
		id := qn.Id

		result = store.GetNews(id)
		time.Sleep(1 * time.Second)

		if token := client.Publish(Config.NewsTopic+"/"+id, 0, false, result); token.Wait() && token.Error() != nil {
			fmt.Printf("%+v\n", token.Error())
		}
		fmt.Println("Published " + Config.NewsTopic + "/" + id)
	})

	return nil
}
