package cbr

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
)

type ICbr struct{}

type CbrCourses struct {
	CharCode string  `json:"char_code"`
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
}

//Источник для курсов фиатных валют: <http://www.cbr.ru/scripts/XML_daily.asp>
func (c ICbr) GetCbr() []CbrCourses {
	doc, err := xmlquery.LoadURL("http://www.cbr.ru/scripts/XML_daily.asp")
	if err != nil {
		fmt.Println(err)
	}

	var cbr []CbrCourses
	for _, v := range xmlquery.Find(doc, "//ValCurs/Valute") {
		value, _ := strconv.ParseFloat(strings.Replace(v.SelectElement("Value").InnerText(), ",", ".", -1), 64)
		nominal, _ := strconv.ParseFloat(strings.Replace(v.SelectElement("Nominal").InnerText(), ",", ".", -1), 64)
		cbr = append(cbr, CbrCourses{
			CharCode: v.SelectElement("CharCode").InnerText(),
			Name:     v.SelectElement("Name").InnerText(),
			Value:    value / nominal,
		})
	}
	return cbr
}
