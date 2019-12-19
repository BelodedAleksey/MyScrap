package hibiny

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

//New struct
type New struct {
	Header   string
	Data     string
	Content  string
	ImageURL string
}

//GetNews func
func GetNews() []New {
	var news []New
	n := New{}
	var url = "https://www.hibiny.com/news"
	c := colly.NewCollector()

	c.OnHTML(`*`, func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr(`href`), `/news/archive`) && e.Text != "" {
			n.Header = e.Text
			p := e.DOM.Parent()
			for i := 0; i < 7; i++ {
				p = p.Parent()
			}
			n.Data = p.Find(`td.p10`).Text()
			n.Content = p.Find(`td.p`).Text()
			news = append(news, n)
		}
		if strings.Contains(e.Attr(`src`), `images/news`) {
			n.ImageURL = `hibiny` + e.Attr(`src`)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error during http request: %s", err)
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Printf("Error visiting %s: %s", url, err)
	}
	c.Wait()
	return news
}
