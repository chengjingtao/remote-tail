package alauda

import "net/http"
import "fmt"
import "os"
import "io/ioutil"
import "encoding/json"

type AlaudaClient struct {
	endpoint string
	client   *http.Client
	token    string
}

func NewAlaudaClient(endpoint, token string) *AlaudaClient {
	return &AlaudaClient{
		endpoint: endpoint,
		token:    token,
		client:   http.DefaultClient,
	}
}

func (client *AlaudaClient) GetInstanceCount(namespace string, app, service string) int {
	path := fmt.Sprintf("%s/v1/services/%s/%s?application=%s", client.endpoint, namespace, service, app)
	fmt.Println("request path -> ", path)
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	request.Header.Add("Authorization", client.token)

	resp, err := client.client.Do(request)
	if err != nil {
		fmt.Println("response error :", err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("response error : not 200, but ", resp.StatusCode)
		os.Exit(1)
	}

	var data interface{}
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error ", err)
		os.Exit(1)
	}
	err = json.Unmarshal(bts, &data)
	if err != nil {
		fmt.Println("json body error ", err)
		os.Exit(1)
	}

	instances, err := NewJSONObjE(data).Get("instances").AsArray()
	if err != nil {
		fmt.Println("json data unexpected: ", err)
		os.Exit(1)
	}
	return len(instances)
}
