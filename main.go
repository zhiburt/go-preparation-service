package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"./service"
)

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"./controllers"
// 	"./service"

// 	"github.com/go-yaml/yaml"
// 	"github.com/gorilla/mux"
// )

// type Config struct {
// 	Address string `yaml:"address"` //server address
// 	Port    string `yaml:"port"`    //server port
// 	Name    string `yaml:"name"`    //server name :)
// }

// var conf Config

// func main() {
// 	router := mux.NewRouter()
// 	router.HandleFunc("/preparations", controllers.AllPreparationsController).
// 		Methods(http.MethodGet)
// 	router.HandleFunc("/preparations/{name}", controllers.PreparationController).
// 		Methods(http.MethodGet)
// 	router.HandleFunc("/preparations/type/{type_name}", controllers.PreparationsByTypeController).
// 		Methods(http.MethodGet)
// 	mServer := http.Server{
// 		Addr:              conf.Address + ":" + conf.Port,
// 		ReadTimeout:       time.Duration(10 * time.Millisecond),
// 		ReadHeaderTimeout: time.Duration(10 * time.Millisecond),
// 		WriteTimeout:      time.Duration(30 * time.Millisecond),
// 		Handler:           router,
// 	}
// 	log.Fatal(mServer.ListenAndServe())
// }

// func init() {
// 	err := service.GetDataFromResources()
// 	if err != nil {
// 		fmt.Println("data not parse")
// 		os.Exit(1)
// 	}
// }

// func init() {
// 	if _, err := conf.GetConfig(); err != nil {
// 		fmt.Println("config not parse")
// 		os.Exit(1)
// 	}
// 	fmt.Println(conf)
// }

// func (c *Config) GetConfig() (*Config, error) {
// 	const confFile = "conf.yaml"
// 	yamlFile, err := ioutil.ReadFile(confFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
// 	}
// 	err = yaml.Unmarshal(yamlFile, c)
// 	if err != nil {
// 		return nil, fmt.Errorf("Unmarshal: %v", err)
// 	}

// 	return c, nil
// }

func main() {
	p, err := service.FindPicture("василиск")
	if err != nil {
		panic(err)
	}
	fmt.Println(p.URL)
}

//init DB-preparations entrys
// func init() {
// 	preps := service.GetPreparationsFromInit()

// 	for _, el := range preps {
// 		b, err := db.InsertIntoPreparations(&db.Preparation{
// 			Name:             el.Name,
// 			Description:      el.Description,
// 			ActiveIngredient: el.ActiveIngredient})
// 		if err != nil {
// 			panic(err)
// 		}
// 		if !b {
// 			fmt.Println("EXISTS ", el.Name)
// 		}
// 	}
// }
