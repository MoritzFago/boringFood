package main

import (
	//	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"strings"
	//	"regexp"
	//	"time"
	"github.com/jinzhu/now"
	"github.com/soh335/ical"
	"gopkg.in/xmlpath.v2"
	"os"
	"regexp"
	"time"
)

func main() {
	makeFood()
}
func makeFood() {
	//xpath für Beilagen: /html/body/div[@id='menu']/div/div[@class='detail']/p/text()
	//xpath für Gerichte: /html/body/div[@id='menu']/div/div[@class='detail']/p/*/text()
	html := getfood()
	path := xmlpath.MustCompile("/html/body/div[@id='menu']/div/div[@class='detail']/p/text()")
	//fmt.Println(html)
	strong := regexp.MustCompile("<strong>")
	strong2 := regexp.MustCompile("</strong>")

	//	html = strong.replaceAllString(html," strong ")
	//	html = strong2.replaceAllString(html," Strong ")
	html = strong.ReplaceAllString(html, "")
	html = strong2.ReplaceAllString(html, " ")
	root, err := xmlpath.ParseHTML(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	c := ical.NewBasicVCalendar()
	c.PRODID = "golang food"
	components := []ical.VComponent{}
	const layout = time.StampNano
	var AddToGetToRightDay = 0
	it := path.Iter(root)
	for it.Next() {
		s := it.Node().String()
		//if s == "" {
		//	break
		//}
		var e ical.VEvent
		e.UID = string(time.Now().Format(layout))
		e.DTSTAMP = time.Now()
		//t2 := now.AddDate(0, 0, 7).Format(layout)
		Day := now.Monday().AddDate(0, 0, AddToGetToRightDay)
		hours, _ := time.ParseDuration("11h")
		hoursah, _ := time.ParseDuration("0.5h")
		DayWTime := Day.Add(hours)
		e.DTSTART = DayWTime
		e.DTEND = DayWTime.Add(hoursah)
		e.SUMMARY = s
		if AddToGetToRightDay == 4 {
			AddToGetToRightDay = 0

		} else {
			AddToGetToRightDay++
		}

		components = append(components, &e)
		//		fmt.Println(string(e.DTSTART.Format(layout)) + ":")
		//		fmt.Println(s)
	}
	c.VComponent = components
	fmt.Print("content-type: text/plain; charset=utf-8\r\n")
	fmt.Print("\r\n")
	c.Encode(os.Stdout)
}
func getfood() string {
	client := http.Client{}
	url := "http://gast.aramark.de/rbz-kiel/menu/web/wochenmenu_de.php"
	resp, err := client.Get(url)
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
