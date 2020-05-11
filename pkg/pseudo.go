package pkg

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"

)

type Cell struct {
	bin                      int
	rnn                      int
	taxpayerOrganization     int
	taxpayerName             int
	ownerName                int
	ownerIin                 int
	ownerRnn                 int
	courtDecision            int
	illegalActivityStartDate int
}

type PseudoCompany struct {
	bin                      string
	rnn                      string
	taxpayerOrganization     string
	taxpayerName             string
	ownerName                string
	ownerIin                 string
	ownerRnn                 string
	courtDecision            string
	illegalActivityStartDate string
}

func (p PseudoCompany) toString() string {
	var id string

	if p.bin != "" {
		id = "\"_id\": \"" + p.bin + "\""
	}
	return "{ \"index\": {" + id + "}} \n" +
		"{ \"bin\":\"" + p.bin + "\"" +
		", \"rnn\":\"" + p.rnn + "\"" +
		", \"taxpayer_organization\":\"" + p.taxpayerOrganization + "\"" +
		", \"taxpayer_name\":\"" + p.taxpayerName + "\"" +
		", \"owner_name\":\"" + p.ownerName + "\"" +
		", \"owner_iin\":\"" + p.ownerIin + "\"" +
		", \"owner_rnn\":\"" + p.ownerRnn + "\"" +
		", \"court_decision\":\"" + p.courtDecision + "\"" +
		", \"illegal_activity_start_date\":\"" + p.illegalActivityStartDate + "\"" +
		"}\n"
}

func ParseAndSendToES(TaxInfoDescription string, f *excelize.File) error {
	cell := Cell{1, 2, 3, 4, 5,
		6, 7, 8, 9}

	replacer := strings.NewReplacer(
		"\"", "'",
		"\\", "/",
		"\n", "",
		"\n\n", "",
		"\r", "")

	for _, name := range f.GetSheetMap() {
		// Get all the rows in the name
		rows := f.GetRows(name)
		var input strings.Builder
		for i, row := range rows {
			if i < 3 {
				continue
			}
			pseudoCompany := new(PseudoCompany)
			for j, colCell := range row {
				switch j {
				case cell.bin:
					pseudoCompany.bin = replacer.Replace(colCell)
				case cell.rnn:
					pseudoCompany.rnn = replacer.Replace(colCell)
				case cell.taxpayerOrganization:
					pseudoCompany.taxpayerOrganization = replacer.Replace(colCell)
				case cell.taxpayerName:
					pseudoCompany.taxpayerName = replacer.Replace(colCell)
				case cell.ownerName:
					pseudoCompany.ownerName = replacer.Replace(colCell)
				case cell.ownerIin:
					pseudoCompany.ownerIin = replacer.Replace(colCell)
				case cell.ownerRnn:
					pseudoCompany.ownerRnn = replacer.Replace(colCell)
				case cell.courtDecision:
					pseudoCompany.courtDecision = replacer.Replace(colCell)
				case cell.illegalActivityStartDate:
					pseudoCompany.illegalActivityStartDate = replacer.Replace(colCell)

				}
			}
			if pseudoCompany.bin != "" {
				input.WriteString(pseudoCompany.toString())
			}
			if i%20000 == 0 {
				if errorT := sendPost(TaxInfoDescription, input.String()); errorT != nil {
					return errorT
				}
				input.Reset()
			}
		}
		if input.Len() != 0 {
			if errorT := sendPost(TaxInfoDescription, input.String()); errorT != nil {
				return errorT
			}
		}
	}
	return nil
}

func sendPost(TaxInfoDescription string, query string) error {
	data := []byte(query)
	r := bytes.NewReader(data)
	resp, err := http.Post("http://localhost:9200/pseudo_company/companies/_bulk", "application/json", r)
	if err != nil {
		fmt.Println("Could not send the data to elastic search " + TaxInfoDescription)
		fmt.Println(err)
		return err
	}
	fmt.Println(TaxInfoDescription + " " + resp.Status)
	return nil
}
