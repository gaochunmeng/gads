package main

import (
	"fmt"
	//"reflect"

	"io/ioutil"
	"net/http"

	"encoding/csv"
	gads "github.com/Getsidecar/gads/v201705"
	"github.com/Getsidecar/sidecar-go-utils/config"
	"os"
	// "strings"
)

func getReport(auth *gads.Auth, headers []string) (collection []map[string]string) {
	fmt.Println("getting report with auth:", auth)
	rds := gads.NewReportDownloadService(auth)

	rd := gads.ReportDefinition{
		ReportName: "ReportNameGoesHere",
		ReportType: "SHOPPING_PERFORMANCE_REPORT",
		DateRangeType: "YESTERDAY",
		DownloadFormat: "CSV",
		Selector: gads.Selector{
			Fields: headers,
			// Predicates: []gads.Predicate{
			// 	{"Status", "EQUALS", []string{"ENABLED"}},
			// 	{"AdvertisingChannelType", "EQUALS", []string{"SHOPPING"}},
			// },
			// Paging: &paging,
		},
	}

	collection, _ = rds.Get(rd)
	return collection
}

func getAWQLResult(auth *gads.Auth, query string) ([]map[string]string) {
	rds := gads.NewReportDownloadService(auth)
	report, err := rds.AWQL(query, "CSV")
	if err != nil {
		fmt.Println("Error in AWQL Query: ", err)
		return nil
	}

	return report
}

func main() {
	authConfig, err := gads.NewCredentialsFromFile("config.json")

	if err != nil {
		panic(err)
	}

	c := &http.Client{}
	configClient := config.ConfigStoreClient{
		HttpClient: c,
		ReadAll:    ioutil.ReadAll,
		BaseUrl:    "https://config.sidecartechnologies.com:4000",
		Username:   "root",
		Password:   "tkw2yWejYMqXm9y",
	}

	clientConfigs, err := configClient.GetClients()

	if err != nil {
		panic(err)
	}
	f, _ := os.Create("test.csv")
	w := csv.NewWriter(f)
	w.Comma = '\t'
	defer f.Close()
	for _, client := range clientConfigs {
		if client.Status != "active" {
			//fmt.Printf("Skipping %s due to inactive flag...\n", client.Shortname)
			continue
		}
		if client.Shortname != "moosejaw" {
			//fmt.Printf("Skipping %s to focus on moosejaw...\n", client.Shortname)
			continue
		}
		fmt.Printf("Running %s...\n", client.Shortname)
		authConfig.Auth.CustomerId = client.Accounts.Adwords.AccountId

		headers := []string{
			"AccountDescriptiveName",
			"CampaignId",
			"AdGroupId",
			"Cost",
			"Clicks",
			"Impressions",
			"Conversions",
			"ConversionValue",
			"OfferId",
			"ExternalCustomerId",
			"Date",
			"AdGroupName",
			"Device",
		}


		//
		//query := `SELECT AccountDescriptiveName,
		//				 CampaignId,
		//				 AdGroupId,
		//				 Cost,
		//			 	 Clicks,
		//				 Impressions,
		//				 Conversions,
		//				 ConversionValue,
		//				 OfferId,
		//				 ExternalCustomerId,
		//				 Date,
		//				 AdGroupName,
		//				 Device
		//		FROM SHOPPING_PERFORMANCE_REPORT
		//		DURING YESTERDAY`


		// For using AWQL
		//report := getAWQLResult(&authConfig.Auth, query)

		// For using Report Download Service
		report := getReport(&authConfig.Auth, headers)

		file, _ := os.Create("result.csv")
		writer := csv.NewWriter(file)
		defer writer.Flush()

		var returnHeaders []string
		for _, value := range report[0:1] {
			for key, _ := range value {
				returnHeaders = append(returnHeaders, key)
			}

			writer.Write(returnHeaders)
		}


		for _, line := range report[0:100] {
			var lineList []string
			for _, header := range returnHeaders {
				lineList = append(lineList, line[header])
			}

			writer.Write(lineList)
		}
	}
	w.Flush()
}
