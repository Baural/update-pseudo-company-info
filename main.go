package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"update_pseudo_company_info/pkg"


	"github.com/360EntSecGroup-Skylar/excelize"
	"gopkg.in/robfig/cron.v2"

)

type TaxInfo struct {
	TaxInfoDescription string
	url                string
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c := schedule()
	c.Start()
	<-make(chan int)
}

func schedule() (c *cron.Cron) {
	c = cron.New()
	_, err := c.AddFunc("30 7 * * *", load)
	if err != nil {
		panic(err)
	}
	return
}

func load() {
	var downloads = []TaxInfo{
		{"Pseudo Company", "http://kgd.gov.kz/mobile_api/services/taxpayers_unreliable_exportexcel/PSEUDO_COMPANY/KZ_ALL/fileName/list_PSEUDO_COMPANY_KZ_ALL.xlsx"},
	}
	var answers = []string{}
	for download := range downloads {
		f := DownloadTaxinfo(downloads[download].url)
		if f == nil {
			answers = append(answers, "Could not read the bad taxpayer information "+downloads[download].TaxInfoDescription)
			continue
		}
		if errorT := pkg.ParseAndSendToES(downloads[download].TaxInfoDescription, f); errorT != nil {
			answers = append(answers, "Could not parse or send to ES, the bad taxpayer information "+downloads[download].TaxInfoDescription)
		} else {
			answers = append(answers, "Have succesfully downloaded the bad taxpayer information "+downloads[download].TaxInfoDescription)
		}
	}
	for answer := range answers {
		fmt.Println(answers[answer])
	}

}

func DownloadTaxinfo(url string) *excelize.File {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	f, err := excelize.OpenReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	return f
}
