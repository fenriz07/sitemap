package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	link "github.com/fenriz07/link/students/fenriz"
	"github.com/fenriz07/sitemap/helpers"
	//"github.com/fenriz07/link"
)

func main() {

	start := time.Now()

	urlFlag := flag.String("url", "https://jerseypedia.org/", "domain to call")

	maxDepth := 2

	flag.Parse()

	pages := bfs(*urlFlag, maxDepth)

	printXML(pages)

	elapsed := time.Since(start)
	log.Printf("Time nano %s", elapsed)
}

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url
}

type Url struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

func printXML(pages []string) {

	urls := make([]Url, 0, len(pages))

	for _, url := range pages {

		urls = append(urls, Url{Loc: url})
	}

	urlset := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}

	out, _ := xml.MarshalIndent(urlset, " ", "  ")

	createFileXML(out)
}

func createFileXML(outputXml []byte) {

	fo, err := os.Create("sitemap.xml")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(fo)

	if _, err := w.Write(outputXml); err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if err = w.Flush(); err != nil {
		panic(err)
	}
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

	//Se obtiene el resultado

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
