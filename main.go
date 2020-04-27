package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/davecgh/go-spew/spew"
	link "github.com/fenriz07/link/students/fenriz"
	"github.com/fenriz07/sitemap/helpers"
	//"github.com/fenriz07/link"
)

func main() {

	urlFlag := flag.String("url", "https://www.mzzo.com/", "domain to call")

	maxDepth := 4

	flag.Parse()

	//index := site{Href: *urlFlag}

	pages := bfs(*urlFlag, maxDepth)

	//index.setChildren(pages)

	spew.Dump(pages)

}

//Algoritmo BFS
func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})

	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i < maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})

		for url, _ := range q {
			/* con la linea # podemos comprobar si la llave existe en el mapa.
			Si dicha llave existe quiere decir que la url fue analizada */

			if _, ok := seen[url]; ok {
				continue
			}

			/* Se le asigna el link que se va a analizar, para que no pueda ser analizado
			en un futuro */
			seen[url] = struct{}{}
			links := get(url)

			/*Se prepara nq con los valores obtenidos que posteriormente se analizaran*/
			for _, link := range links {
				nq[link] = struct{}{}
			}
		}
	}

	ret := make([]string, 0, len(seen))

	for url, _ := range seen {
		ret = append(ret, url)
	}

	return ret
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
