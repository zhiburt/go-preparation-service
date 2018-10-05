package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Picture struct {
	Name string
	URL  string
}

func FindPicture(name string) (*Picture, error) {
	client := &http.Client{Timeout: time.Second * 50}
	url := fmt.Sprintf(`http://search.tut.by/?is=1&page=0&ru=1&isize=&query=%s`, name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Find Picture,with parametr <%s> error: %v", name, err)
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)

	d, _ := goquery.NewDocumentFromReader(res.Body)

	block := d.Find("div .serp-list .serp-list_type_search")
	if block == nil {
		return nil, fmt.Errorf("Picture not found")
	}

	aImg := block.Find("a .serp-item__link")
	if aImg != nil {
		return nil, fmt.Errorf("Something WRONG A item not found")
	}
	img := aImg.Find("img .serp-item__thumb")
	src, _ := img.Attr("src")
	src = "http:" + src

	p := &Picture{
		Name: name,
		URL:  src,
	}

	return p, nil
}
