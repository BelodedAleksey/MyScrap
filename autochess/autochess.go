package autochess

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

//Player for checkPlayer
type Player struct {
	Name          string
	ID            string
	Rank          string
	RankPoint     string
	MaxRank       string
	MaxRankPoint  string
	Candy         string
	Level         string
	LevelPoint    string
	MatchesPlayed string
	AveragePlace  string
	Top1          string
	Top3          string
	Rounds        string
	Games         []Game
}

//Game struct
type Game struct {
	Place  string
	Data   string
	Lobby  string
	MMR    string
	Time   string
	Rounds string
	VPN    string
	Gold   string
	Health string
	Cost   string
}

//PlayerT for Top500
type PlayerT struct {
	Name    string
	ID      string
	Rank    string
	MMR     string
	IconURL string
}

//GetTop500 func
func GetTop500() ([]PlayerT, error) {
	var players []PlayerT
	player := PlayerT{}
	urlMain := "https://auto-chess.ru/mobile-top1000/"
	cMain := colly.NewCollector()

	cMain.OnHTML(`tr.main`, func(e *colly.HTMLElement) {
		player.Name = e.ChildText(`span`)
		player.Rank = e.ChildText(`td:not([class])`)
		player.MMR = e.ChildText(`.rank-1`)
		s := strings.Split(e.ChildAttr(`a`, `href`), `/`)
		player.ID = s[len(s)-2]
		player.IconURL = `auto-chess.ru` + e.ChildAttr(`img`, `src`)
		players = append(players, player)
	})

	cMain.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error during http request: %s", err)
	})

	err := cMain.Visit(urlMain)
	if err != nil {
		fmt.Printf("Error visiting %s: %s", urlMain, err)
		return nil, err
	}
	cMain.Wait()
	return players, nil
}

//GetPlayersByName func
func GetPlayersByName(name string) ([]PlayerT, error) {
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
		result    string
		players   []PlayerT
		outerhtml string
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
		chromedp.OuterHTML(`*[class^="heroes-list"]`, &outerhtml),
		//chromedp.Evaluate(jsGetText(`.hero-flex`), &names),
		//chromedp.AttributesAll(`a[href^="/check-player/"][target="blank"]`, &ids, chromedp.ByQueryAll),
	)

	if err != nil {
		fmt.Println("Ошибка ChromeDP: ", err)
		return nil, err
	}
	//Names
	reg, err := regexp.Compile(`<br>(.*?)<br>`)
	if err != nil {
		fmt.Println("Error regexp: ", err)
	}
	names := reg.FindAllStringSubmatch(outerhtml, -1)
	//Ids
	reg, err = regexp.Compile(`check-player/(.*?)/"`)
	if err != nil {
		fmt.Println("Error regexp: ", err)
	}
	ids := reg.FindAllStringSubmatch(outerhtml, -1)
	//Ranks
	reg, err = regexp.Compile(`"common">(.*?)</span>`)
	if err != nil {
		fmt.Println("Error regexp: ", err)
	}
	ranks := reg.FindAllStringSubmatch(outerhtml, -1)
	//Icons
	reg, err = regexp.Compile(`image:url\((.*?)\)">`)
	if err != nil {
		fmt.Println("Error regexp: ", err)
	}
	iconUrls := reg.FindAllStringSubmatch(outerhtml, -1)
	fmt.Println("RESULT: ", result)
	for i, n := range names {
		players = append(players,
			PlayerT{ID: ids[i][1], Name: n[1], Rank: ranks[i][1],
				IconURL: "auto-chess.ru" + iconUrls[i][1]})
	}
	return players, nil
}

//GetPlayerByID func
func GetPlayerByID(id string) (Player, error) {
	player := Player{}
	var games []Game
	urlPlayer := "https://auto-chess.ru/check-player/" + id
	cPlayer := colly.NewCollector()
	cPlayer.OnHTML(`.article__content *`, func(e *colly.HTMLElement) {
		if e.Attr(`id`) == `info-block` {
			player.Name = e.ChildText(`h3:first-of-type`)
			player.Rank = e.ChildText(`h3:nth-of-type(2)`)
			player.ID = strings.TrimPrefix(e.ChildText(`div:first-of-type`), "ID: ")
			player.Candy = e.ChildText(`div:nth-of-type(2)`)
		}
		if e.Attr(`id`) == `rank-level` {
			text, err := e.DOM.Html()
			if err != nil {
				fmt.Println("Error HTML: ", err)
			}
			text = strings.TrimSuffix(text, "</span>")
			text = strings.TrimPrefix(text, "Уровень ")
			txt := strings.Split(text, "<br/><span>")
			player.Level = txt[0]
			player.LevelPoint = txt[1]
		}
		if strings.Contains(e.Attr(`class`), `rank-big`) {
			switch e.ChildText(`h3`) {
			case "Ранг игрока":
				player.RankPoint = e.ChildText(`span`)
			case "Макс. ранг":
				player.MaxRankPoint = e.ChildText(`span`)
			case "Матчей сыграно":
				s := "Матчей сыграно"
				player.MatchesPlayed = strings.Trim(e.Text, s)
			case "Среднее место":
				s := "Среднее место"
				player.AveragePlace = strings.Trim(e.Text, s)
			case "Топ-1":
				s := "Топ-1"
				player.Top1 = strings.Trim(e.Text, s)
			case "Топ-3":
				s := "Топ-3"
				player.Top3 = strings.Trim(e.Text, s)
			case "Раундов в матчах":
				s := "Раундов в матчах"
				player.Rounds = strings.Trim(e.Text, s)
			default:
			}
		}

		//История игр

		if strings.Contains(e.Attr(`id`), `match`) {
			game := Game{}
			game.Place = e.ChildText(`div[class^="match-place"]`)
			mInfo := e.DOM.Find(`.match-info`).First()
			for i := 0; i < 9; i++ {
				s := mInfo.Find(`span`).Text()
				switch s {
				case "Дата":
					game.Data = strings.Trim(mInfo.Text(), s)
					mInfo = mInfo.Next()
				case "Лобби":
					game.Lobby = strings.Trim(mInfo.Text(), s)
					mInfo = mInfo.Next()
				case "MMR":
					game.MMR = mInfo.Find(`div`).Text()
					mInfo = mInfo.Next()
				case "Время":
					game.Time = strings.Trim(mInfo.Text(), s)
					mInfo = mInfo.Next()
				case "Раунды":
					game.Rounds = strings.Trim(mInfo.Text(), s)
					mInfo = mInfo.Next()
				case "В/П/Н":
					game.VPN = strings.Trim(mInfo.Text(), s)
					mInfo = mInfo.Next()
				case "Золото":
					game.Gold = strings.Trim(mInfo.Text(), s)
					mInfo = mInfo.Next()
				case "Здоровье":
					game.Health = strings.Trim(mInfo.Text(), s)
					mInfo = mInfo.Next()
				case "Цена сборки":
					game.Cost = strings.Trim(mInfo.Text(), s)
					mInfo = mInfo.Next()
				}
			}
			games = append(games, game)
		}
	})

	cPlayer.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error during http request: %s", err)
	})

	err := cPlayer.Visit(urlPlayer)
	if err != nil {
		fmt.Printf("Error visiting %s: %s", urlPlayer, err)
		return Player{}, err
	}
	cPlayer.Wait()
	player.Games = games
	return player, nil
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
