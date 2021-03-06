package main

import (
	"github.com/zhangzhoujian/istio-demo/echo/internal/schema"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type comResponse struct {
	Meta     schema.Meta
	Companies []*schema.Company `json:"elements"`
}
var (
	prodV1Count int
	prodV2Count int
)

func getCompanies(xHeader string) {
	url := "http://192.168.99.100:31380/company/api/v1/companies"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("x-request-platform", xHeader)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	er := &comResponse{}
	err = json.Unmarshal(body, er)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range er.Companies {
		for _, p := range c.Products {
			fmt.Printf("{ x-request-platform: [%s] } ==> product version [%s], pod name [%s]\n",xHeader, p.Meta.Version, p.Meta.PodName)
			if p.Meta.Version == "1.0.0" {
				prodV1Count = prodV1Count + 1
			}
			if p.Meta.Version == "2.0.0" {
				prodV2Count = prodV2Count + 1
			}	
		}
	}
}

func main() {
	argNum := len(os.Args)
	// fmt.Printf("the num of input is %d\n",argNum)

	// fmt.Printf("%v", os.Args)
	if argNum < 2 {
		for i := 0; i < 10; i++ {
			getCompanies("")
		}
		fmt.Printf("prod of company count v1: %d, v2: %d\n", prodV1Count, prodV2Count)
		return
	} 

	for i := 0; i < 10; i++ {
		getCompanies(os.Args[1])
	}
	fmt.Printf("prod of company count v1: %d, v2: %d\n", prodV1Count, prodV2Count)
}