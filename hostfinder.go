//Script to find the hostname  name from the public IP

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"html/template"
	"log"
	"net/http"
)

var input string
var ipconfigs string
var ipconfigsv2 string
var ipAddress string
var host_name_ip string



type PageVariables struct {
	Host_name         string

}

var id_tag map[string]interface{}

func parseMap(aMap map[string]interface{}) {
	for key, val := range aMap {

		switch concreteVal := val.(type) {

		case map[string]interface{}:

			id_tag = concreteVal


			parseMap(val.(map[string]interface{}))

		case []interface{}:

			parseArray(val.([]interface{}))

		default:

			if key == "ipAddress" {

				if concreteVal == ipconfigsv2 {


          value_id, ok := id_tag["id"]
            if ok {

            host_val_original := value_id.(string)




            host_val := strings.Split(host_val_original, "/")

            if host_val[7] == "applicationGateways" {


							app_gw:=" APP-Gateway and hostname is ::   "
							host_name_ip=app_gw+host_val[8]

            } else {

							vm_host_name_ip:=strings.Split(host_val[8], "ni_")


							vmh:=" Virtul Machine and hostname is ::   "

							host_name_ip=vmh+vm_host_name_ip[1]

            }

              } else {

                desc_ip()

                }

				}
			}

		}

	}
}

func parseArray(anArray []interface{}) {


	for i, val := range anArray {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			parseMap(val.(map[string]interface{}))
		case []interface{}:
			_ = i
			parseArray(val.([]interface{}))
		default:
			_ = concreteVal

		}
	}
}

func ip_segragator(ip_desc_details string) {

	input = ip_desc_details

	m := map[string]interface{}{}

	err := json.Unmarshal([]byte(input), &m)

	if err != nil {
		panic(err)
	}

	parseMap(m)

}

func ip_describe(actokenip string) {

	url := "https://management.azure.com/subscriptions/xxxxxxx/providers/Microsoft.Network/publicIPAddresses?api-version=2018-11-01"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	req.Header.Add("referer", "https://docs.microsoft.com/en-us/rest/api/virtualnetwork/publicipaddress(preview)/listall")
	req.Header.Add("origin", "https://docs.microsoft.com")
	var bearertoken = "Bearer " + actokenip

	req.Header.Add("authorization", bearertoken)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "xxxxxxxxxxxxxxx")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	ip_desc_details := string(body)

	ip_segragator(ip_desc_details)

}

func desc_ip() {

	url := "https://login.microsoftonline.com/xxxxxxxxxxxxxxxxxx/oauth2/token"

	payload := strings.NewReader("------WebKitFormBoundary7MA4YWasassxkTrasZu0gW\r\nContent-Disposition: form-data; name=\"grant_type\"\r\n\r\nclient_credentials\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"client_id\"\r\n\r\n8da51433-47d4-4ccf-a24a-25fe51b7f9a4\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"client_secret\"\r\n\r\njBF*b=CCPg0=l.1Bw066*.Os*EPvWOw-\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"resource\"\r\n\r\nhttps://management.azure.com/\r\n------WebKitFormBoundary7MA4YWxkTrZu0asasdgW--")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YwdwdWxkTrZwdwdu0gW")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "xxxxxxxxxxxxxxxxxxxxxxxxx")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	auth_token_response := string(body)
	sec := map[string]interface{}{}
	if err := json.Unmarshal([]byte(auth_token_response), &sec); err != nil {
		panic(err)
	}
	value, ok := sec["access_token"]
	if ok {
		actoken, _ := value.(string)
		ip_describe(actoken)

	} else {
		fmt.Println("access token not not found")
	}

}



func login(w http.ResponseWriter, r *http.Request) {



  if r.Method == "GET" {
      t, _ := template.ParseFiles("login.html")
      t.Execute(w, nil)




  } else {
      r.ParseForm()

      ipAddress=r.Form["ip_address"][0]
			ipconfigsv2=ipAddress

			desc_ip()


      HomePageVars := PageVariables{ //store the date and time in a struct
        Host_name: host_name_ip,
      }

      t, err := template.ParseFiles("outputpage.html") //parse the html file homepage.html
      if err != nil { // if there is an error
        log.Print("template parsing error: ", err) // log it
      }


      err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
      if err != nil { // if there is an error
        log.Print("template executing error: ", err) //log it
      }


  }



}

func main() {
	host_name_ip= "Invalid Entry"

	http.HandleFunc("/",login)
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
			log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("")


}
