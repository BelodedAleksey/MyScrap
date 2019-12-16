package kinobilety

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
)

var url = "https://api.kinobilety.net/api/getFilms?cityId=199&marketId=931" //city=Апатиты market=Полярный

type Film struct {
	ID         string
	Name       string
	Genre      string
	Country    string
	URLTrailer string
	URLPoster  string
}

//Get func
func Get() {
	films := []Film{}
	currentFilm := Film{}

	cSeances := colly.NewCollector()
	cSeances.OnHTML("div[data-date]", func(e *colly.HTMLElement) {
		fmt.Println("Сеанс: ", e.Attr("data-date"))
		e.DOM.Find("a[href]").Each(func(i int, gc *goquery.Selection) {
			fmt.Println("Ссылка: ", gc.AttrOr("href", ""))
			fmt.Println("Время: ", gc.Children().Text())
		})
	})

	cMain := colly.NewCollector()
	cMain.OnHTML(`div[class^="film"]`, func(e *colly.HTMLElement) {
		switch e.Attr("class") {
		case "film_name":
			currentFilm.Name = e.Text
			fmt.Println("Фильм: ", e.Text)
		case "film_genre":
			currentFilm.Genre = e.Text
			fmt.Println("Возраст/Жанр: ", e.Text)
		case "film_country":
			currentFilm.Country = e.Text
			fmt.Println("Страна: ", e.Text)
		case "film_trailer":
			currentFilm.URLTrailer = e.Attr("data-url")
			fmt.Println("Трейлер: ", e.Attr("data-url"))
			films = append(films, currentFilm)
		case "film_poster":
			currentFilm.URLPoster = e.ChildAttr("img", "src")
			fmt.Println("Постер: ", e.ChildAttr("img", "src"))
		default:
			if id := e.Attr("id"); id != "" {
				id = strings.TrimPrefix(id, "film_")
				currentFilm.ID = id
				fmt.Println("ID: ", id)
				urlSeance := "https://api.kinobilety.net/api/getSchedule?filmId=" + id +
					"&cityId=199&marketId=931"
				cSeances.Visit(urlSeance)
			}

		}
	})
	cMain.Visit(url)
	cMain.Wait()

}
