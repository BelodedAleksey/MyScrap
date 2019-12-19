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

//District struct
type District struct {
	Name string
	ID   string
}

//Clinic struct
type Clinic struct {
	District District
	ID       string
	Name     string
	Address  string
}

//Specialist struct
type Specialist struct {
	Clinic Clinic
	ID     string
	Name   string
}

//Doctor struct
type Doctor struct {
	Specialist Specialist
	ID         string
	Name       string
}

func chromeSearch(url string, tasks chromedp.Tasks) {
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

	actions := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#footer`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("VISIBLE!!!")
			return nil
		}),
		//chromedp.Click(`//li[text()="`+district+`"]`, chromedp.NodeVisible),}
	}
	actions = append(actions, tasks...)

	err := chromedp.Run(ctx, actions)
	if err != nil {
		log.Fatal(err)
	}
}

//GetDistricts func
func GetDistricts() ([]District, error) {
	var districts []District
	districtURL :=
		`https://reg.polarmed.ru/schedule/`
	cMain := colly.NewCollector()
	cMain.OnHTML(`li[data-id]`, func(e *colly.HTMLElement) {
		if id := e.Attr("data-id"); id != "" {
			districts = append(districts, District{ID: id, Name: e.Text})
		}
	})

	cMain.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error during http request: %s", err)
	})

	err := cMain.Visit(districtURL)
	if err != nil {
		fmt.Printf("Error visiting %s: %s", districtURL, err)
		return nil, err
	}
	cMain.Wait()
	return districts, nil
}

//GetClinics func
func (d District) GetClinics() ([]Clinic, error) {
	var clinics []Clinic

	url := `https://reg.polarmed.ru/schedule/#%5B%7B%22district%22%3A` +
		d.ID + `%7D%5D`

	// navigate to a page, wait for an element, click
	var names, addresses []string
	var attrs []map[string]string

	tasks := chromedp.Tasks{
		chromedp.WaitVisible(`span.clinic_name`),
		chromedp.AttributesAll(`#clinic_list > li`, &attrs, chromedp.ByQueryAll),
		chromedp.Evaluate(jsGetText(`span.clinic_name`), &names),
		chromedp.Evaluate(jsGetText(`span.clinic_address`), &addresses),
	}

	chromeSearch(url, tasks)

	for i, n := range names {
		clinic := Clinic{}
		clinic.District = d
		clinic.ID = attrs[i][`data-id`]
		clinic.Name = n
		clinic.Address = addresses[i]
		clinics = append(clinics, clinic)
	}
	return clinics, nil
}

//GetSpecialists func
func (c Clinic) GetSpecialists() ([]Specialist, error) {
	var specialists []Specialist

	url := `https://reg.polarmed.ru/schedule/#%5B%7B%22district%22%3A` +
		c.District.ID + `%7D%2C%7B%22clinic%22%3A` +
		c.ID + `%7D%5D`

	// navigate to a page, wait for an element, click
	var names []string
	var attrs []map[string]string

	tasks := chromedp.Tasks{
		chromedp.WaitVisible(`#speciality_list`),
		chromedp.AttributesAll(`#speciality_list > li`, &attrs, chromedp.ByQueryAll),
		chromedp.Evaluate(jsGetText(`#speciality_list > li`), &names),
	}

	chromeSearch(url, tasks)

	for i, n := range names {
		spec := Specialist{}
		spec.Clinic = c
		spec.Name = n
		spec.ID = attrs[i]["data-id"]
		specialists = append(specialists, spec)
	}
	return specialists, nil
}

//GetDoctors func
func (s Specialist) GetDoctors() ([]Doctor, error) {
	var doctors []Doctor
	url := `https://reg.polarmed.ru/schedule/#%5B%7B%22district%22%3A` +
		s.Clinic.District.ID + `%7D%2C%7B%22clinic%22%3A` +
		s.Clinic.ID + `%7D%2C%7B%22speciality%22%3A` +
		s.ID + `%7D%5D`
	var names []string
	var attrs []map[string]string
	tasks := chromedp.Tasks{
		chromedp.WaitVisible(``),
		chromedp.AttributesAll(``, &attrs, chromedp.ByQueryAll),
		chromedp.Evaluate(jsGetText(``), &names),
	}
	chromeSearch(url, tasks)
	for i, n := range names {
		doc := Doctor{}
		doc.Specialist = s
		doc.Name = n
		doc.ID = attrs[i]["data-id"]
		doctors = append(doctors, doc)
	}
	return doctors, nil
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
