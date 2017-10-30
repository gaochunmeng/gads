package main

import (
	"fmt"
	"reflect"

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
	fmt.Println("res type:", reflect.TypeOf(collection))
	return collection
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
		// cs := gads.NewCampaignService(&authConfig.Auth)

		// paging := gads.Paging{
		// 	Offset: int64(0),
		// 	Limit:  int64(1000),
		// }

		// sets, _, err := cs.Get(
		// 	gads.Selector{
		// 		Fields: []string{
		// 			"Id",
		// 			"Name",
		// 			"BudgetId",
		// 			"Amount",
		// 		},
		// 		Predicates: []gads.Predicate{
		// 			{"Status", "EQUALS", []string{"ENABLED"}},
		// 			{"AdvertisingChannelType", "EQUALS", []string{"SHOPPING"}},
		// 		},
		// 		Paging: &paging,
		// 	},
		// )

		// if err != nil && !strings.Contains(err.Error(), "Authentication") && !strings.Contains(err.Error(), "Authorization") {
		// 	fmt.Println(err)
		// 	continue
		// }

		//fmt.Printf("sets: %+v \n", sets)

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

		report := getReport(&authConfig.Auth, headers)
		// fmt.Println("report:", report

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


		for _, line := range report[0:10] {
			var lineList []string
			for _, header := range returnHeaders {
				lineList = append(lineList, line[header])
			}

			writer.Write(lineList)
		}


		//w.WriteStructHeader

		// myFile, _ := os.Create("myTest.csv")
		// myWriter := csv.NewWriter(os.Stdout)
		// myWriter.Comma = '\t'
		// myWriter.WriteStructAll(report)

		// for _, set := range sets {
		// 	//w.WriteStructAll(set)
		// 	w.Write([]string{client.SiteID, client.Shortname, fmt.Sprintf("%d", set.Id), set.Name, fmt.Sprintf("%d", set.BudgetId), fmt.Sprintf("%d", set.BudgetAmount)})
		// }
	}
	w.Flush()
	//fmt.Printf("%#v\n", campaignMap)
	//
	//for _, campaigns := range campaignMap {
	//	for _, campaign := range campaigns {
	//		fmt.Printf("%s: %d\n", campaign.Name, campaign.Budget.Amount/1000000)
	//	}
	//}

	//sharedSetService := gads.NewSharedSetService(&config.Auth)
	//op := gads.SharedSetOperation{
	//	Operator: "ADD",
	//	Operand: gads.SharedSet{
	//		Name: "Zach's dumb test list",
	//		Type: "NEGATIVE_KEYWORDS",
	//	},
	//}
	//err = sharedSetService.Mutate([]gads.SharedSetOperation{op})
	//fmt.Println(err)

	//sharedCriterionService := gads.NewSharedCriterionService(&config.Auth)
	//
	//sets, _, err := sharedCriterionService.Get(selector)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//ops := []gads.SharedCriterionOperation{}
	//for _, criterion := range sets {
	//	criterion.SharedSetId = 1534457185
	//	op := gads.SharedCriterionOperation{
	//		Operator: "ADD",
	//		Operand:  criterion,
	//	}
	//	ops = append(ops, op)
	//}
	//err = sharedCriterionService.Mutate(ops)
}
