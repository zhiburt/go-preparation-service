package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"uni/coorse/db"
	// "../db"

	"github.com/PuerkitoBio/goquery"
)

type Preparation struct {
	Name             string
	Price            float64
	Description      string
	ActiveIngredient string
}

func GetPreparationsFromInit() []*Preparation {
	//One file obrabotka
	preps := make([]*Preparation, 0)
	f, err := os.Open(file1)
	var scaner *bufio.Scanner
	if err != nil {
		data, err := GetDataOne()
		if err != nil {
			panic("Get dataOne error: " + err.Error())
		}
		scaner = bufio.NewScanner(strings.NewReader(data))
	} else {
		defer f.Close()
		scaner = bufio.NewScanner(f)
	}

	scaner.Split(bufio.ScanLines)
	scaner.Scan() // we don't need first line it's comment
	//two preparats in data string line
	for scaner.Scan() {
		data := strings.SplitN(scaner.Text(), SEPARATOR, -1)
		price, err := strconv.ParseFloat(data[2], 64)
		if err != nil {
			panic("parse float from data error: " + err.Error())
		}
		prep := &Preparation{
			Name:             data[1],
			Price:            price,
			ActiveIngredient: data[0],
		}

		preps = addPrep(preps, prep)

		price, err = strconv.ParseFloat(data[4], 64)
		if err != nil {
			panic("parse float from data error: " + err.Error())
		}
		prep2 := &Preparation{
			Name:             data[3],
			Price:            price,
			ActiveIngredient: data[0],
		}

		preps = addPrep(preps, prep2)
	}

	//Two file obrabotka
	f1, err := os.Open(file2)
	if err != nil {
		data, err := GetDataTwo()
		if err != nil {
			panic("Get dataOne error: " + err.Error())
		}
		scaner = bufio.NewScanner(strings.NewReader(data))
	} else {
		defer f1.Close()
		scaner = bufio.NewScanner(f1)
	}

	scaner.Split(bufio.ScanLines)
	scaner.Scan() // we don't need first line it's comment
	//two preparats in data string line
	for scaner.Scan() {
		data := strings.SplitN(scaner.Text(), SEPARATOR, -1)
		price, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			panic("parse float from data error: " + err.Error())
		}
		prep := &Preparation{
			Name:  data[0],
			Price: price,
		}

		preps = addPrep(preps, prep)

		price, err = strconv.ParseFloat(data[3], 64)
		if err != nil {
			panic("parse float from data error: " + err.Error())
		}
		prep2 := &Preparation{
			Name:  data[2],
			Price: price,
		}

		preps = addPrep(preps, prep2)
	}

	return preps
}

//logic by lower price need make
func addPrep(preps []*Preparation, p *Preparation) []*Preparation {
	for _, el := range preps {
		if (*el).Name == p.Name {
			log.Println("value %q already containts in preparations", p.Name)
			return preps
		}
	}

	preps = append(preps, p)

	return preps
}

const (
	file1 = "initData/site1_data.txt"
	file2 = "initData/site2_data.txt"
)

//need move these
func GetDataFromResources() error {
	f, err := OpenOrCreateFile(file1)
	if err != nil {
		return fmt.Errorf("open file1 err %s", err)
	}
	defer f.Close()

	data, err := GetDataOne()
	if err != nil {
		return fmt.Errorf("GetDataOne: %s", err)
	}
	_, err = f.WriteString(data)
	log.Println(err)

	f1, err := OpenOrCreateFile(file2)
	if err != nil {
		return fmt.Errorf("open file2 err %s", err)
	}
	defer f1.Close()

	data, err = GetDataTwo()
	if err != nil {
		return fmt.Errorf("GetDataTwo: %s", err)
	}
	_, err = f1.WriteString(data)
	log.Println(err)

	return nil
}

func OpenOrCreateFile(name string) (*os.File, error) {
	f, err := os.Open(name)
	if err != nil {
		log.Println(err)
		f, err = os.Create(name)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return f, err
}

const SEPARATOR = "><"

func GetDataOne() (string, error) {
	const url = "http://инструкция-от-таблетки.рф/дешевые_аналоги_дорогих_лекарств"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}
	table := doc.Find("table")
	buf := new(bytes.Buffer)
	line := new(bytes.Buffer)
	table.Find("tr").Each(func(i int, element *goquery.Selection) {
		// if i == 0 {
		// 	return
		// }
		element.Find("td").Each(func(_ int, s *goquery.Selection) {
			line.WriteString(s.Text() + SEPARATOR)
		})

		buf.WriteString(line.String() + "\n")
		line.Reset()
	})

	return buf.String(), nil
}

func GetDataTwo() (string, error) {
	const url = "http://лечим-грибок.рф/generic-list/"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}
	table := doc.Find("table")
	buf := new(bytes.Buffer)
	line := new(bytes.Buffer)
	table.Find("tr").Each(func(i int, element *goquery.Selection) {
		// if i == 0 {
		// 	return
		// }
		element.Find("td").Each(func(_ int, s *goquery.Selection) {
			line.WriteString(s.Text() + SEPARATOR)
		})

		buf.WriteString(line.String() + "\n")
		line.Reset()
	})

	return buf.String(), nil
}

//support db fill
func getRandomFromSliceStr(sl []string, s *rand.Source) string {
	rand := rand.New(*s)

	return sl[rand.Intn(len(sl)-1)]
}

func getMedicalPreparatTypes() []string {
	var medicals []string
	f, err := os.Open("./initData/medical_types.txt")
	if err != nil {
		panic(fmt.Errorf("getTypes error: %v", err))
	}
	scaner := bufio.NewScanner(f)
	scaner.Split(bufio.ScanLines)
	for scaner.Scan() {
		medicals = append(medicals, scaner.Text())
	}

	return medicals
}

func getMedicalPreparatImage() []string {
	var medicals []string
	f, err := os.Open("./initData/medical_images.txt")
	if err != nil {
		panic(fmt.Errorf("getTypes error: %v", err))
	}
	scaner := bufio.NewScanner(f)
	scaner.Split(bufio.ScanLines)
	for scaner.Scan() {
		medicals = append(medicals, scaner.Text())
	}

	return medicals
}

const MEDICAL_SEPARATOR = "|"

func getMedicalPreparatDescription() []string {
	f, err := os.Open("./initData/medical_description.txt")
	if err != nil {
		panic(fmt.Errorf("getTypes error: %v", err))
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(fmt.Errorf("getMedicalPreparatDescription error: %v", err))
	}
	descs := strings.SplitN(string(b), MEDICAL_SEPARATOR, -1)
	descs = append(descs, "")
	descs = append(descs, "")
	descs = append(descs, "")

	return descs
}

//init DB-preparations entrys
func init() {
	// GetDataFromResources()
	//if table db.Preparations is empty
	if db.GetAllPreparations() == nil {
		GetDataFromResources()
		preps := GetPreparationsFromInit()
		//rand type
		rand := rand.NewSource(time.Now().Unix())
		for _, el := range preps {
			err := db.InsertIntoPreparations(&db.Preparation{
				Name:             el.Name,
				Description:      getRandomFromSliceStr(getMedicalPreparatDescription(), &rand),
				ActiveIngredient: el.ActiveIngredient,
				Type:             getRandomFromSliceStr(getMedicalPreparatTypes(), &rand),
				ImageURL:         getRandomFromSliceStr(getMedicalPreparatImage(), &rand),
			})
			if err != nil {
				panic(err)
			}
		}
	}
}
