package main

//TODO RENAME ROUTING

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	// _ "github.com/lib/pq"

	"uni/coorse/controllers"

	_ "uni/coorse/service"

	"github.com/go-yaml/yaml"
	"github.com/gorilla/mux"

	_ "net/http/pprof"
)

//URLs
const (
	//Preparations
	AllPreparationsURL         = "/preparations"
	AllPreparationsVersion2URL = "/preparations/all"
	FindPreparationByNameURL   = "/preparations/find/byName"
	FindPreparationByIdURL     = "/preparations/find/byId"
	NewPreparationURL          = "/preparations/new"
	UpdatePreparationURL       = "/preparations/update"
	DeletePreparationURL       = "/preparations/delete"

	//Suppliers
	AllSuppliersURL                 = "/suppliers"
	AllSuppliersVersion2URL         = "/suppliers/all"
	FindSupplierByAddressAndNameURL = "/suppliers/find/byAddressAndName"
	FindSuppliersByCompanyURL       = "/suppliers/find/byCompany"
	FindSupplierByIdURL             = "/suppliers/find/byId"
	NewSupplierURL                  = "/suppliers/new"
	UpdateSupplierURL               = "/suppliers/update"
	DeleteSupplierURL               = "/suppliers/delete"

	//Goods
	AllGoodsURL              = "/goods"
	AllGoodsBySupplierURL    = "/goods/suppliers"
	AllGoodsByPreparationURL = "/goods/preparations"
	NewGoodURL               = "/goods/new"
	DeteteGoodURL            = "/goods/delBySupplierete"
)

type Config struct {
	Address string `yaml:"address"` //server address
	Port    string `yaml:"port"`    //server port
	Name    string `yaml:"name"`    //server name :)
}

var conf Config

func main() {
	mServer := NewServer()
	log.Fatal(mServer.ListenAndServe())
}

func NewServer() *http.Server {
	router := NewHandler()

	mServer := &http.Server{
		Addr:              conf.Address + ":" + conf.Port,
		ReadTimeout:       time.Duration(10 * time.Millisecond),
		ReadHeaderTimeout: time.Duration(10 * time.Millisecond),
		WriteTimeout:      time.Duration(30 * time.Millisecond),
		Handler:           router,
	}

	return mServer
}

func NewHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc(AllPreparationsURL, controllers.AllPreparationsController).
		Methods(http.MethodGet)
	router.HandleFunc(AllPreparationsVersion2URL, controllers.AllPreparationsController).
		Methods(http.MethodGet)
	router.HandleFunc(FindPreparationByNameURL, controllers.FindPreparationByNameController).
		Methods(http.MethodGet)
	router.HandleFunc(FindPreparationByIdURL, controllers.FindPreparationByIdController).
		Methods(http.MethodGet)
	router.HandleFunc(NewPreparationURL, controllers.InsertPreparationController).
		Methods(http.MethodPut)
	router.HandleFunc(UpdatePreparationURL, controllers.UpdatePreparationController).
		Methods(http.MethodPost)
	router.HandleFunc(DeletePreparationURL, controllers.DeletePreparationController).
		Methods(http.MethodDelete)

	router.HandleFunc(AllSuppliersURL, controllers.AllSuppliersController).
		Methods(http.MethodGet)
	router.HandleFunc(AllSuppliersVersion2URL, controllers.AllSuppliersController).
		Methods(http.MethodGet)
	router.HandleFunc(FindSupplierByAddressAndNameURL, controllers.FindSupplierByNameController).
		Methods(http.MethodGet)
	router.HandleFunc(FindSuppliersByCompanyURL, controllers.FindSuppliersByCompanyController).
		Methods(http.MethodGet)
	router.HandleFunc(FindSupplierByIdURL, controllers.FindSupplierByIdController).
		Methods(http.MethodGet)
	router.HandleFunc(NewSupplierURL, controllers.InsertSupplierController).
		Methods(http.MethodPut)
	router.HandleFunc(UpdateSupplierURL, controllers.UpdateSupplierController).
		Methods(http.MethodPost)
	router.HandleFunc(DeleteSupplierURL, controllers.DeleteSupplierController).
		Methods(http.MethodDelete)

	router.HandleFunc(AllGoodsURL, controllers.AllGoodsController).
		Methods(http.MethodGet)
	router.HandleFunc(AllGoodsBySupplierURL, controllers.AllSuppliersGoodsController).
		Methods(http.MethodGet)
	router.HandleFunc(AllGoodsByPreparationURL, controllers.AllPreparationsGoodsController).
		Methods(http.MethodGet)
	router.HandleFunc(NewGoodURL, controllers.InsertGoodController).
		Methods(http.MethodPut)
	router.HandleFunc(DeteteGoodURL, controllers.DeleteGoodController).
		Methods(http.MethodDelete)

	return router
}

func init() {
	// err := service.GetDataFromResources()
	// if err != nil {
	// 	fmt.Println("data not parse")
	// 	os.Exit(1)
	// }
}

func init() {
	if _, err := conf.GetConfig(); err != nil {
		fmt.Println("config not parse")
		os.Exit(1)
	}
	fmt.Println(conf)
}

func (c *Config) GetConfig() (*Config, error) {
	const confFile = "conf.yaml"
	yamlFile, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	return c, nil
}

// 	p, err := service.FindPicture("василиск")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(p.URL)
