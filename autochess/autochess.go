package autochess

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

type Player struct {
	Name          string
	ID            string
	Rank          string
	RankPoint     string
	MaxRank       string
	MaxRankPoint  string
	Candy         string
	Level         string
	MatchesPlayed string
	AveragePlace  string
	Top1          string
	Top3          string
	Rounds        string
}

//GetTop500 func
func GetTop500() {
	urlMain := "https://auto-chess.ru/mobile-top1000/"
	cMain := colly.NewCollector()
	cMain.OnHTML(`tr.main`, func(e *colly.HTMLElement) {
		fmt.Println("Игрок: ", e.ChildText(`span`))
		fmt.Println("Ранг: ", e.ChildText(`td:not([class])`))
		fmt.Println("ММР: ", e.ChildText(`.rank-1`))
		s := strings.Split(e.ChildAttr(`a`, `href`), `/`)
		fmt.Println("ID: ", s[len(s)-2])
		fmt.Println("Иконка: ", `auto-chess.ru`+e.ChildAttr(`img`, `src`))
	})
	cMain.Visit(urlMain)
	cMain.Wait()
}

//GetPlayersByName func
func GetPlayersByName(name string) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), //Визуальное отображение
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var (
		result string
		names  []string
		ids    []map[string]string
	)
	urlCheck := "https://auto-chess.ru/check-player/"
	//Поиск по имени
	err := chromedp.Run(ctx,
		chromedp.Navigate(urlCheck),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`.footer__notes`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("VISIBLE!!!")
			return nil
		}),
		chromedp.SendKeys(`//input[@name="playername"]`, name),
		chromedp.Submit(`input[value="Найти игроков!"]`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("SUBMITTED!!!")
			return nil
		}),
		chromedp.WaitVisible(`.search-results`),
		chromedp.Text(`.search-results`, &result, chromedp.NodeVisible),
		chromedp.Evaluate(jsGetText(`.hero-flex`), &names),
		chromedp.AttributesAll(`a[href^="/check-player/"]`, &ids, chromedp.ByQueryAll),
	)

	if err != nil {
		fmt.Println("Ошибка ChromeDP: ", err)
	}
	fmt.Println("RESULT: ", result)
	for i, n := range names {
		fmt.Println("Name: ", n)
		s := strings.Split(ids[i]["href"], `/`)
		fmt.Println("ID: ", s[len(s)-2])
	}

}

//GetPlayerByID func
func GetPlayerByID(id string) {
	//player := Player{}
	urlPlayer := "https://auto-chess.ru/check-player/" + id
	cPlayer := colly.NewCollector()
	cPlayer.OnHTML(`.article__content *`, func(e *colly.HTMLElement) {
		if e.Attr(`id`) == `info-block` {
			fmt.Println("Игрок: ", e.ChildText(`h3:first-of-type`))
			fmt.Println("Ранг: ", e.ChildText(`h3:nth-of-type(2)`))
			fmt.Println("ID: ", e.ChildText(`div:first-of-type`))
			fmt.Println("Конфеты: ", e.ChildText(`div:nth-of-type(2)`))
		}
		if e.Attr(`id`) == `rank-level` {
			text, err := e.DOM.Html()
			if err != nil {
				fmt.Println("Error HTML: ", err)
			}
			text = strings.TrimSuffix(text, "</span>")
			txt := strings.Split(text, "<br/><span>")
			fmt.Println("Уровень: ", txt[0])
			fmt.Println("Очков: ", txt[1])
		}
		if strings.Contains(e.Attr(`class`), `rank-big`) {
			switch e.ChildText(`h3`) {
			case "Ранг игрока":
				s := "Ранг игрока"
				fmt.Println(s+": ", strings.Trim(e.Text, s))
			case "Макс. ранг":
				s := "Макс. ранг"
				fmt.Println(s+": ", strings.Trim(e.Text, s))
			case "Матчей сыграно":
				s := "Матчей сыграно"
				fmt.Println(s+": ", strings.Trim(e.Text, s))
			case "Среднее место":
				s := "Среднее место"
				fmt.Println(s+": ", strings.Trim(e.Text, s))
			case "Топ-1":
				s := "Топ-1"
				fmt.Println(s+": ", strings.Trim(e.Text, s))
			case "Топ-3":
				s := "Топ-3"
				fmt.Println(s+": ", strings.Trim(e.Text, s))
			case "Раундов в матчах":
				s := "Раундов в матчах"
				fmt.Println(s+": ", strings.Trim(e.Text, s))
			default:
			}
		}

		//История игр
		if strings.Contains(e.Attr(`id`), `match`) {
			fmt.Println("Место: ", e.ChildText(`.match-place`))
			mInfo := e.DOM.Find(`.match-info`).First()
			for i := 0; i < 9; i++ {
				//Дата**Лобби**MMR**Время**Раунды**В/П/Н**Золото**Здоровье**Цена сборки**
				s := mInfo.Children().Text()
				fmt.Println(s+": ", strings.Trim(mInfo.Text(), s))
				mInfo = mInfo.Next()
			}
		}

	})

	cPlayer.Visit(urlPlayer)
	cPlayer.Wait()
}

//Get All Text of Elements
func jsGetText(sel string) (js string) {
	const funcJS = `function getText(sel) {
				var text = [];
				var elements = document.body.querySelectorAll(sel);

				for(var i = 0; i < elements.length; i++) {
					var current = elements[i];
					if(current.textContent.replace(/ |\n/g,'') !== '') {
					// Check the element is not empty
						text.push(current.textContent + ',');
					}
				}
				return text
			 };`

	invokeFuncJS := `var a = getText('` + sel + `'); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}
