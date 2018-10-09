package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	// _ "uni/coorse/db"
	//_ "github.com/lib/pq"
)

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

var (
	client = &http.Client{Timeout: time.Second}
)

type Case struct {
	Method string // GET по-умолчанию в http.NewRequest если передали пустую строку
	Path   string
	Query  string
	Status int
	Result interface{}
}

//CaseResponce
type CR map[string]interface{}

func TestPreparations(t *testing.T) {
	log = &logrus.Logger{}
	log.Out = ioutil.Discard

	cases := []*Case{
		//FindPreparationByIdController
		&Case{
			"", FindPreparationByIdURL, "id=10", http.StatusOK, CR{
				"data": CR{
					"Id":               10,
					"Name":             "Декспантенол",
					"Description":      "\nYou are buying 100 x wipes. These are typically used for cleaning the skin before injections but are also perfect for cleaning electronic components as well as a multiple of other uses. The wipes are individually wrapped and impregnated with 70% isopropyl alcohol which will kill most common bacteria.",
					"Type":             "Inhalers",
					"ActiveIngredient": "Декспантенол",
					"ImageURL":         "https://s3.amazonaws.com/uifaces/faces/twitter/aviddayentonbay/128.jpg",
				},
			},
		},
		&Case{
			"", FindPreparationByIdURL, "id=1000000", http.StatusOK, CR{
				"data": nil,
			},
		},
		&Case{
			"", FindPreparationByIdURL, "id=-1", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Id: -1 does not validate as numeric",
			},
		},
		//FindPreparationByNameController
		&Case{
			"", FindPreparationByNameURL, "name=Асприн-кардио", http.StatusOK, CR{
				"data": CR{
					"Id":               5,
					"Name":             "Асприн-кардио",
					"Description":      "Unlike traditional bulb ear syringes that can damage the ear if inserted too far, the unique Aculife Ear Wax Removal Syringe features a flared tip design to prevent over insertion while effectively cleaning the ears and preventing ear wax buildup. The unique tip directs fluid to the ear canal walls using a tri-stream directional jet ensuring an even safer and more effective alternative to direct flow syringes. The ear wax removal syringe is used to effectively dislodge stubborn ear wax and other debris. Exit portals allow for drainage ensuring no pressure build-up occurs and any wax or debris is effectively drained safely away. For ear wax blockages it is advisable to soften the wax prior to syringing with olive oil. This unique product is designed for home use and comes complete with full user instructions, an improved syringe with a larger 50ml capacity \u0026 10 Salt Sachets. Not to be used with children under 3 years of age.",
					"Type":             "Drops",
					"ActiveIngredient": "Ацетилсалициловая кислота",
					"ImageURL":         "https://s3.amazonaws.com/uifaces/faces/twitter/bpartridge/128.jpg",
				},
			},
		},
		&Case{
			"", FindPreparationByNameURL, "name=invalid_name", http.StatusOK, CR{
				"data": nil,
			},
		},
		&Case{
			"", FindPreparationByNameURL, "", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Name: non zero value required",
			},
		},
		//InsertPreparationController
		&Case{
			"PUT", NewPreparationURL, "name=TEST_PREPARATION_2", http.StatusOK, CR{
				"status": "OK",
			},
		},
		&Case{
			"PUT", NewPreparationURL, "name=TEST_PREPARATION_2", http.StatusInternalServerError, CR{
				"error": `error -> Name TEST_PREPARATION_2 exists in DB`,
			},
		},
		&Case{
			"PUT", NewPreparationURL, "name=", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Name: non zero value required",
			},
		},
		//UpdatePreparationController
		&Case{
			"POST", UpdatePreparationURL, "id=276&name=RENAME_PREPARATION", http.StatusOK, CR{
				"status": "OK",
			},
		},
		&Case{
			"POST", UpdatePreparationURL, "id=0&name=RedMedic", http.StatusInternalServerError, CR{
				"error": "preparation with name RedMedic or id 0 doesn't exists in db",
			},
		},
		&Case{
			"POST", UpdatePreparationURL, "id=1&name=", http.StatusOK, CR{
				"status": "OK",
			},
		},
		&Case{
			"POST", UpdatePreparationURL, "name=RENAME_PREPARATION&description=new_description", http.StatusOK, CR{
				"status": "OK",
			},
		},
		&Case{
			"POST", UpdatePreparationURL, "", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> name or id must be not empty",
			},
		},
		&Case{
			"POST", UpdatePreparationURL, "id=10000000&name=test", http.StatusInternalServerError, CR{
				"error": "preparation with name test or id 10000000 doesn't exists in db",
			},
		},
		//DeletePreparationController
		&Case{
			"DELETE", DeletePreparationURL, "name=TEST_PREPARATION_2", http.StatusOK, CR{
				"status": "OK",
			},
		},
		&Case{
			"DETELE", DeletePreparationURL, "id=0", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> name or id must be not empty",
			},
		},
		&Case{
			"DELETE", DeletePreparationURL, "id=10000", http.StatusInternalServerError, CR{
				"error": "preparation with name  or id 10000 doesn't exists in db",
			},
		},
	}
	m := http.NewServeMux()
	m.HandleFunc(FindPreparationByNameURL, FindPreparationByNameController)
	m.HandleFunc(FindPreparationByIdURL, FindPreparationByIdController)
	m.HandleFunc(NewPreparationURL, InsertPreparationController)
	m.HandleFunc(DeletePreparationURL, DeletePreparationController)
	m.HandleFunc(UpdatePreparationURL, UpdatePreparationController)

	RunTest(t, m, cases)

	cases = []*Case{
		//AllPreparationsController
		&Case{
			"", AllPreparationsURL, "", http.StatusOK, CR{
				"data": struct{}{},
			},
		},
		//ivalid
		&Case{
			"", AllPreparationsVersion2URL, "", http.StatusOK, CR{
				"data": struct{}{},
			},
		},
	}

	m.HandleFunc(AllPreparationsURL, AllPreparationsController)
	m.HandleFunc(AllPreparationsVersion2URL, AllPreparationsController)

	RunTestCheckOnlyStatus(t, m, cases)
}

func TestSuppliers(t *testing.T) {
	log.Out = ioutil.Discard

	cases := []*Case{
		//FindSupplierByIdController
		&Case{
			"", FindSupplierByIdURL, "id=1", http.StatusOK, CR{
				"data": CR{
					"Id":          1,
					"Name":        "Microsoft Washington PP",
					"Company":     "Microsoft",
					"Address":     "Washington 200B 3c 10.2",
					"Geolocation": "",
					"Description": "",
				},
			},
		},
		//FindSupplierByIdController
		&Case{
			"", FindSupplierByIdURL, "id=1000000", http.StatusOK, CR{
				"data": nil,
			},
		},
		&Case{
			"", FindSupplierByIdURL, "id=-1", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Id: -1 does not validate as numeric",
			},
		},
		//InsertPreparationController
		&Case{
			"PUT", NewSupplierURL, "name=TEST_NAME&address=TEST_ADDRESS&company=TEST_COMPANY", http.StatusOK, CR{
				"status": "OK",
			},
		},
		&Case{
			"PUT", NewSupplierURL, "name=TEST_NAME&address=TEST_ADDRESS&company=TEST_COMPANY", http.StatusInternalServerError, CR{
				"error": `db: supplier exists in`,
			},
		},
		&Case{
			"PUT", NewSupplierURL, "", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Name: non zero value required;Address: non zero value required",
			},
		},
		&Case{
			"PUT", NewSupplierURL, "name=2323", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Address: non zero value required",
			},
		},
		//FindSupplierByNameController
		&Case{
			"", FindSupplierByAddressAndNameURL, "name=Microsoft%20Washington%20PP&address=Washington%20200B%203c%2010.2", http.StatusOK, CR{
				"data": CR{
					"Id":          1,
					"Name":        "Microsoft Washington PP",
					"Company":     "Microsoft",
					"Address":     "Washington 200B 3c 10.2",
					"Geolocation": "",
					"Description": "",
				},
			},
		},
		&Case{
			"", FindSupplierByAddressAndNameURL, "name=invalid_name_&address=invalid_address_", http.StatusOK, CR{
				"data": nil,
			},
		},
		&Case{
			"", FindSupplierByAddressAndNameURL, "", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Address: non zero value required;Name: non zero value required",
			},
		},
		//FindSuppliersByCompany
		&Case{
			"", FindSuppliersByCompanyURL, "company=under_bed", http.StatusOK, CR{
				"data": []CR{
					CR{
						"Id":          8,
						"Name":        "under_bed",
						"Company":     "under_bed",
						"Address":     "under_address_bed",
						"Geolocation": "under_bed",
						"Description": "",
					},
				},
			},
		},
		&Case{
			"", FindSuppliersByCompanyURL, "", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Company: non zero value required",
			},
		},
		&Case{
			"", FindSuppliersByCompanyURL, "company=invalid_name_company", http.StatusOK, CR{
				"data": nil,
			},
		},
		//UpdatePreparationController
		&Case{
			"POST", UpdateSupplierURL, "id=8&name=under_bed&company=under_bed&address=under_address_bed&geolocation=under_bed", http.StatusOK, CR{
				"status": "OK",
			},
		},
		&Case{
			"POST", UpdateSupplierURL, "id=10000&name=valid&company=valid&address=valid", http.StatusInternalServerError, CR{
				"error": "db: supplier doesn't exists in",
			},
		},
		&Case{
			"POST", UpdateSupplierURL, "", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Id: non zero value required;Name: non zero value required;Company: non zero value required;Address: non zero value required",
			},
		},
		&Case{
			"POST", UpdateSupplierURL, "id=10000000", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> Name: non zero value required;Company: non zero value required;Address: non zero value required",
			},
		},
		//DeletePreparationController
		&Case{
			"DELETE", DeleteSupplierURL, "name=TEST_NAME&address=TEST_ADDRESS", http.StatusOK, CR{
				"status": "OK",
			},
		},
		&Case{
			"DETELE", DeleteSupplierURL, "id=0", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> id or name + address must be not empty",
			},
		},
		&Case{
			"DETELE", DeleteSupplierURL, "name=test_name", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> id or name + address must be not empty",
			},
		},
		&Case{
			"DETELE", DeleteSupplierURL, "address=test_address", http.StatusInternalServerError, CR{
				"error": "InvalidRequest -> id or name + address must be not empty",
			},
		},
		&Case{
			"DELETE", DeleteSupplierURL, "id=100000", http.StatusInternalServerError, CR{
				"error": "db: supplier doesn't exists in",
			},
		},
	}
	m := http.NewServeMux()
	m.HandleFunc(FindSupplierByIdURL, FindSupplierByIdController)
	m.HandleFunc(FindSupplierByAddressAndNameURL, FindSupplierByNameController)
	m.HandleFunc(FindSuppliersByCompanyURL, FindSuppliersByCompanyController)
	m.HandleFunc(NewSupplierURL, InsertSupplierController)
	m.HandleFunc(DeleteSupplierURL, DeleteSupplierController)
	m.HandleFunc(UpdateSupplierURL, UpdateSupplierController)

	RunTest(t, m, cases)

	cases = []*Case{
		//AllSuppliersController
		&Case{
			"", AllSuppliersURL, "", http.StatusOK, nil,
		},
		&Case{
			"", AllSuppliersVersion2URL, "", http.StatusOK, nil,
		},
	}

	m.HandleFunc(AllSuppliersURL, AllSuppliersController)
	m.HandleFunc(AllSuppliersVersion2URL, AllSuppliersController)

	RunTestCheckOnlyStatus(t, m, cases)
}

func RunTest(t *testing.T, m http.Handler, cases []*Case) {
	ts := httptest.NewServer(m)
	for idx, item := range cases {
		var (
			err      error
			result   interface{}
			expected interface{}
			req      *http.Request
		)

		caseName := fmt.Sprintf("case %d: [%s] %s %s", idx, item.Method, item.Path, item.Query)

		if item.Method == http.MethodPost || item.Method == http.MethodPut {
			reqBody := strings.NewReader(item.Query)
			req, err = http.NewRequest(item.Method, ts.URL+item.Path, reqBody)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req, err = http.NewRequest(item.Method, ts.URL+item.Path+"?"+item.Query, nil)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != item.Status {
			t.Errorf("[%s] expected http status %v, got %v", caseName, item.Status, resp.StatusCode)
			continue
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Errorf("[%s] cant unpack json: %v", caseName, err)
			continue
		}

		// reflect.DeepEqual не работает если нам приходят разные типы
		// а там приходят разные типы (string VS interface{}) по сравнению с тем что в ожидаемом результате
		// этот маленький грязный хак конвертит данные сначала в json, а потом обратно в interface - получаем совместимые результаты
		// не используйте это в продакшен-коде - надо явно писать что ожидается интерфейс или использовать другой подход с точным форматом ответа
		data, err := json.Marshal(item.Result)
		json.Unmarshal(data, &expected)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("[%d] results not match\nGot: %#v\nExpected: %#v", idx, result, item.Result)
			continue
		}
	}
}

func RunTestCheckOnlyStatus(t *testing.T, m http.Handler, cases []*Case) {
	ts := httptest.NewServer(m)
	for idx, item := range cases {
		var (
			err error
			req *http.Request
		)

		caseName := fmt.Sprintf("case %d: [%s] %s %s", idx, item.Method, item.Path, item.Query)

		if item.Method == http.MethodPost || item.Method == http.MethodPut {
			reqBody := strings.NewReader(item.Query)
			req, err = http.NewRequest(item.Method, ts.URL+item.Path, reqBody)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req, err = http.NewRequest(item.Method, ts.URL+item.Path+"?"+item.Query, nil)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != item.Status {
			t.Errorf("[%s] expected http status %v, got %v", caseName, item.Status, resp.StatusCode)
			continue
		}
	}
}
