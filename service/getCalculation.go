package service

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type FundoResponse struct {
	BazinCoservative float64 `json:"bazinCoservative"`
	BazinHurled      float64 `json:"bazinHurled"`
	CurrentValue     float64 `json:"currentValue"`
}
type AcaoResponse struct {
	BazinCoservative float64 `json:"bazinCoservative"`
	BazinHurled      float64 `json:"bazinHurled"`
	Graham           float64 `json:"graham"`
	CurrentValue     float64 `json:"currentValue"`
}

func GetBazin(fundo string) (*FundoResponse, error) {
	c := colly.NewCollector()

	var currentValue string
	selectCurrentValue := "#cards-ticker > div._card.cotacao > div._card-body > div > span"

	var lastTdTexts []string

	c.OnHTML("#table-dividends-history", func(e *colly.HTMLElement) {

		e.DOM.Find("tbody tr").Each(func(index int, rowHtml *goquery.Selection) {

			lastTd := rowHtml.Find("td").Last()
			lastTdText := strings.TrimSpace(lastTd.Text())

			lastTdTexts = append(lastTdTexts, lastTdText)
		})
	})

	c.OnHTML(selectCurrentValue, func(e *colly.HTMLElement) {
		currentValue = e.Text
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	url := "https://investidor10.com.br/fiis/" + fundo + "/"
	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	if len(lastTdTexts) > 12 {
		lastTdTexts = lastTdTexts[:12]
	}

	sum := 0.0

	for _, text := range lastTdTexts {

		text = strings.ReplaceAll(text, ",", ".")
		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			fmt.Printf("Error parsing value: %v\n", err)
			return nil, err
		}
		sum += value
	}
	var bazinCoservative float64
	var bazinHurled float64
	bazinCoservative = sum / 0.06
	bazinHurled = sum / 0.1
	fmt.Println(sum)

	currentValue = strings.ReplaceAll(currentValue, ",", ".")
	currentValue = strings.ReplaceAll(currentValue, " ", "")
	currentValue = strings.ReplaceAll(currentValue, "R$", "")

	currentValueFloat, err := covertStringToFloat64(currentValue)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	roundedBazinCoservative := math.Ceil(bazinCoservative*100) / 100
	roundedBazinHurled := math.Ceil(bazinHurled*100) / 100

	fundoResponse := &FundoResponse{
		BazinCoservative: roundedBazinCoservative,
		BazinHurled:      roundedBazinHurled,
		CurrentValue:     currentValueFloat,
	}

	return fundoResponse, nil
}

func GetBazinAndGraham(acao string) (*AcaoResponse, error) {
	c := colly.NewCollector()
	a := colly.NewCollector()

	var vpa, lpa, currentValue, dy string

	selectCurrentValue := "#main-2 > div:nth-child(4) > div > div.pb-3.pb-md-5 > div > div.info.special.w-100.w-md-33.w-lg-20 > div > div:nth-child(1) > strong"
	selectDy := "#cards-ticker > div._card.dy > div._card-body > span"

	c.OnHTML(selectDy, func(e *colly.HTMLElement) {
		dy = e.Text
	})

	selectorVpa := "#indicators-section > div.indicator-today-container > div > div:nth-child(1) > div > div:nth-child(9) > div > div > strong"
	selectLpa := "#indicators-section > div.indicator-today-container > div > div:nth-child(1) > div > div:nth-child(11) > div > div > strong"
	a.OnHTML(selectorVpa, func(e *colly.HTMLElement) {
		vpa = e.Text
	})

	a.OnHTML(selectLpa, func(e *colly.HTMLElement) {
		lpa = e.Text
	})

	a.OnHTML(selectCurrentValue, func(e *colly.HTMLElement) {
		currentValue = e.Text
	})

	a.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	url := "https://investidor10.com.br/acoes/" + acao + "/"
	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	url2 := "https://statusinvest.com.br/acoes/" + acao + "/"
	err = a.Visit(url2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("VPA: %s\n", vpa)
	fmt.Printf("LPA: %s\n", lpa)
	dy = strings.ReplaceAll(dy, ",", ".")
	dy = strings.ReplaceAll(dy, " ", "")
	dy = strings.ReplaceAll(dy, "%", "")
	dy = strings.ReplaceAll(dy, "-", "")
	if dy == "" {
		dy = "0.0"
	}

	vpa = strings.ReplaceAll(vpa, ",", ".")
	lpa = strings.ReplaceAll(lpa, ",", ".")
	currentValue = strings.ReplaceAll(currentValue, ",", ".")

	dyFloat, err := covertStringToFloat64(dy)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	vpaFloat, err := covertStringToFloat64(vpa)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	lpaFloat, err := covertStringToFloat64(lpa)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	currentValueFloat, err := covertStringToFloat64(currentValue)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var bazinCoservative float64
	var bazinHurled float64
	var graham float64
	bazinCoservative = currentValueFloat * (dyFloat / 100) / 0.06
	bazinHurled = currentValueFloat * (dyFloat / 100) / 0.1
	graham = math.Sqrt(22.5 * lpaFloat * vpaFloat)

	roundedBazinCoservative := math.Ceil(bazinCoservative*100) / 100
	roundedBazinHurled := math.Ceil(bazinHurled*100) / 100
	roundedGraham := math.Ceil(graham*100) / 100

	acaoResponse := &AcaoResponse{
		BazinCoservative: roundedBazinCoservative,
		BazinHurled:      roundedBazinHurled,
		Graham:           roundedGraham,
		CurrentValue:     currentValueFloat,
	}

	return acaoResponse, nil
}

func covertStringToFloat64(value string) (float64, error) {

	if strings.HasPrefix(value, "-") {
		negativeValue, err := strconv.ParseFloat(value[1:], 64)
		if err != nil {
			fmt.Println("Error:", err)
			return 0.0, err
		}
		return -negativeValue, nil
	}

	convertedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0.0, err
	}
	return convertedValue, nil
}
