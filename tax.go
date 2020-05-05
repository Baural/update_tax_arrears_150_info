package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Cell struct {
	region                 int
	officeOfTaxEnforcement int
	oteId                  int
	bin                    int
	rnn                    int
	taxpayerOrganizationRu int
	taxpayerOrganizationKz int
	lastNameKz             int
	firstNameKz            int
	middleNameKz           int
	lastNameRu             int
	firstNameRu            int
	middleNameRu           int
	ownerIin               int
	ownerRnn               int
	ownerNameKz            int
	ownerNameRu            int
	economicSector         int
	totalDue               int
	subTotalMain           int
	subTotalFine           int
	subTotalLateFee        int
}

type Tax struct {
	region                 string
	officeOfTaxEnforcement string
	oteId                  string
	bin                    string
	rnn                    string
	taxpayerOrganizationRu string
	taxpayerOrganizationKz string
	lastNameKz             string
	firstNameKz            string
	middleNameKz           string
	lastNameRu             string
	firstNameRu            string
	middleNameRu           string
	ownerIin               string
	ownerRnn               string
	ownerNameKz            string
	ownerNameRu            string
	economicSector         string
	totalDue               string
	subTotalMain           string
	subTotalFine           string
	subTotalLateFee        string
}

func (t Tax) toString() string {
	var id string

	if t.bin != "" {
		id = "\"_id\": \"" + t.bin + "\""
	}
	return "{ \"index\": {" + id + "}} \n" +
		"{ \"region\":\"" + t.region + "\"" +
		", \"office_of_tax_enforcement\":\"" + t.officeOfTaxEnforcement + "\"" +
		", \"ote_id\":\"" + t.oteId + "\"" +
		", \"bin\":\"" + t.bin + "\"" +
		", \"rnn\":\"" + t.rnn + "\"" +
		", \"taxpayer_organization_ru\":\"" + t.taxpayerOrganizationRu + "\"" +
		", \"taxpayer_organization_kz\":\"" + t.taxpayerOrganizationKz + "\"" +
		", \"last_name_kz\":\"" + t.lastNameKz + "\"" +
		", \"first_name_kz\":\"" + t.firstNameKz + "\"" +
		", \"middle_name_kz\":\"" + t.middleNameKz + "\"" +
		", \"last_name_ru\":\"" + t.lastNameRu + "\"" +
		", \"first_name_ru\":\"" + t.firstNameRu + "\"" +
		", \"middle_name_ru\":\"" + t.middleNameRu + "\"" +
		", \"owner_iin\":\"" + t.ownerIin + "\"" +
		", \"owner_rnn\":\"" + t.ownerRnn + "\"" +
		", \"owner_name_kz\":\"" + t.ownerNameKz + "\"" +
		", \"owner_name_ru\":\"" + t.ownerNameRu + "\"" +
		", \"economic_sector\":\"" + t.economicSector + "\"" +
		", \"total_due\":\"" + t.totalDue + "\"" +
		", \"sub_total_main\":\"" + t.subTotalMain + "\"" +
		", \"sub_total_fine\":\"" + t.subTotalFine + "\"" +
		", \"sub_total_late_fee\":\"" + t.subTotalLateFee + "\"" +
		"}\n"
}

func parseAndSendToES(TaxInfoDescription string, f *excelize.File) error {
	cell := Cell{1, 2, 3, 4, 5,
		6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22}

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
			tax := new(Tax)
			for j, colCell := range row {
				switch j {
				case cell.region:
					tax.region = replacer.Replace(colCell)
				case cell.officeOfTaxEnforcement:
					tax.officeOfTaxEnforcement = replacer.Replace(colCell)
				case cell.oteId:
					tax.oteId = replacer.Replace(colCell)
				case cell.bin:
					tax.bin = replacer.Replace(colCell)
				case cell.rnn:
					tax.rnn = replacer.Replace(colCell)
				case cell.taxpayerOrganizationRu:
					tax.taxpayerOrganizationRu = replacer.Replace(colCell)
				case cell.taxpayerOrganizationKz:
					tax.taxpayerOrganizationKz = replacer.Replace(colCell)
				case cell.lastNameKz:
					tax.lastNameKz = replacer.Replace(colCell)
				case cell.firstNameKz:
					tax.firstNameKz = replacer.Replace(colCell)
				case cell.middleNameKz:
					tax.middleNameKz = replacer.Replace(colCell)
				case cell.lastNameRu:
					tax.lastNameRu = replacer.Replace(colCell)
				case cell.firstNameRu:
					tax.firstNameRu = replacer.Replace(colCell)
				case cell.middleNameRu:
					tax.middleNameRu = replacer.Replace(colCell)
				case cell.ownerIin:
					tax.ownerIin = replacer.Replace(colCell)
				case cell.ownerRnn:
					tax.ownerRnn = replacer.Replace(colCell)
				case cell.ownerNameKz:
					tax.ownerNameKz = replacer.Replace(colCell)
				case cell.ownerNameRu:
					tax.ownerNameRu = replacer.Replace(colCell)
				case cell.economicSector:
					tax.economicSector = replacer.Replace(colCell)
				case cell.totalDue:
					tax.totalDue = replacer.Replace(colCell)
				case cell.subTotalMain:
					tax.subTotalFine = replacer.Replace(colCell)
				case cell.subTotalLateFee:
					tax.subTotalLateFee = replacer.Replace(colCell)
				}
			}
			// if tax.bin != "" {
			input.WriteString(tax.toString())
			// }
			if i%10000 == 0 {
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
	resp, err := http.Post("http://localhost:9200/tax_arrears_150/companies/_bulk", "application/json", r)
	if err != nil {
		fmt.Println("Could not send the data to elastic search " + TaxInfoDescription)
		fmt.Println(err)
		return err
	}
	fmt.Println(TaxInfoDescription + " " + resp.Status)
	return nil
}
