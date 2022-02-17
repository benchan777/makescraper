package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

type ScrapedData struct {
	MenuData string
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// data := e.Text
		data := ScrapedData{e.Text}
		reg, _ := regexp.Compile("menu_data = (.*)};") // compile regexp
		rawData := reg.FindString(data.MenuData)       // find string

		// ignore the first and last part of the regex data ("menu_data = " and "};")
		length := len(rawData)
		jsonify, _ := json.Marshal(rawData[12 : length-1])

		// Since the data is already in JSON form, running json.Marshal adds an extra set of quotes causing
		// unnecessary escaping. Remove the extra quotes so proper json can be built
		unquoted, _ := strconv.Unquote(string(jsonify))

		// Write JSON to file
		_ = os.WriteFile("output.json", []byte(unquoted), 0644)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://guerrero.tartine.menu/pickup/")
}
