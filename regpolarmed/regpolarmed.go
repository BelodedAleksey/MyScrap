package regpolarmed

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

//Get func
func Get() {
	cityURL :=
		`https://reg.polarmed.ru/schedule/`
	cMain := colly.NewCollector()
	cMain.OnHTML("*", func(e *colly.HTMLElement) {
		if e.Name == "li" && e.Attr("data-id") != "" {
			id := e.Attr("data-id")
			fmt.Println("ID: ", id)        //id
			fmt.Println("Город: ", e.Text) //Города
		}
	})

	cMain.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	cMain.Visit(cityURL)
	cMain.Wait()

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

	url := `https://reg.polarmed.ru/schedule/`

	// navigate to a page, wait for an element, click
	var names, addresses, specials []string
	var attrs []map[string]string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`#footer`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("VISIBLE!!!")
			return nil
		}),
		chromedp.Click(`li[data-id="10"]`, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Открыты клиники!!!")
			return nil
		}),
		chromedp.WaitVisible(`span.clinic_name`),
		chromedp.AttributesAll(`#clinic_list > li`, &attrs, chromedp.ByQueryAll),
		chromedp.Evaluate(jsGetText(`span.clinic_name`), &names),
		chromedp.Evaluate(jsGetText(`span.clinic_address`), &addresses),

		chromedp.Click(`li[data-id="51"]`, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Открыты специалисты!!!")
			return nil
		}),
		chromedp.WaitVisible(`#speciality_list`),
		chromedp.AttributesAll(`#speciality_list > li`, &attrs, chromedp.ByQueryAll),
		chromedp.Evaluate(jsGetText(`#speciality_list > li`), &specials),
	)

	if err != nil {
		log.Fatal(err)
	}

	for i, n := range names {
		log.Println("ID: ", attrs[i][`data-id`])
		log.Println("Название: ", n)
		log.Println("Адрес: ", addresses[i])
	}
}

//Get All Text of Elements
func jsGetText(sel string) (js string) {
	const funcJS = `function getText(sel) {
				var text = [];
				var elements = document.body.querySelectorAll(sel);

				for(var i = 0; i < elements.length; i++) {
					var current = elements[i];
					if(current.children.length === 0 && current.textContent.replace(/ |\n/g,'') !== '') {
					// Check the element has no children && that it is not empty
						text.push(current.textContent + ',');
					}
				}
				return text
			 };`

	invokeFuncJS := `var a = getText('` + sel + `'); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}
