package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fenriz07/link"
	"github.com/fenriz07/sitemap/helpers"
	//"github.com/fenriz07/link"
)

type site struct {
	Href      string
	childrens []site
}

func (s *site) setChildren(pages []string) {

	childrens := []site{}

	for _, p := range pages {

		siteStruct := site{
			Href: p,
		}

		childrens = append(childrens, siteStruct)
	}

	s.childrens = childrens

}

func main() {

	urlFlag := flag.String("url", "https://gophercises.com/", "domain to call")

	flag.Parse()

	index := site{Href: *urlFlag}

	pages := get(*urlFlag)

	index.setChildren(pages)

	for k, children := range index.childrens {

		pages := get(children.Href)

		index.childrens[k].setChildren(pages)
	}

	spew.Dump(index)

}

func get(urlStr string) []string {

	resp, err := http.Get(urlStr)

	if err != nil {
		helpers.DD(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}

	base := baseUrl.String()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		helpers.DD(err)
	}

	allLinks := link.ParseHtml(string(body))

	pages := filter(base, createPages(*allLinks, base))

	return pages

}

func createPages(links []link.Link, base string) []string {

	var href string
	var allLinks []string

	for _, l := range links {

		href = l.Href

		switch {
		case strings.HasPrefix(href, "/"):
			allLinks = append(allLinks, base+href)
		case strings.HasPrefix(href, "http"):
			allLinks = append(allLinks, href)
		}
	}

	return allLinks
}

func filter(base string, links []string) []string {
	var ret []string

	for _, link := range links {
		if strings.HasPrefix(link, base) {
			ret = append(ret, link)
		}
	}

	return unique(ret)
}

func unique(elements []string) []string {

	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}
