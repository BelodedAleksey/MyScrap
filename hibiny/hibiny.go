package hibiny

import (
	"fmt"

	"github.com/gocolly/colly"
)

var url = "https://www.hibiny.com/"

//Get func
func Get() {

	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Text == "Прогноз погоды" { //текст ссылки
			link := e.DOM.Parent().Parent().Next()     //поднимаемся на 2 уровня выше и берем следущую строку на этом уровне
			txt := link.Children().Find("span").Text() //ищем span, его текст это t°C
			fmt.Println("TEXT: ", txt)
		}

		if e.Text == "Кино в Апатитах" {
			c.Visit(url + e.Attr("href")) //склеиваем ссылки
		}
		//Поднимаемся пока не найдем td с потомком h1, если его текст Афиша то берем текст ссылок
		if e.DOM.ParentsUntil("td > h1").Find("h1").Text() == "Кино в Апатитах - Кинотеатр Полярный" {
			fmt.Println(e.Text)
		}

	})
	c.Visit(url)
	c.Wait()
}
