package smoke

import (
	"io/ioutil"
	"encoding/json"
	graphite "github.com/marpaia/graphite-golang"
	"net/http"
	"log"
	"strconv"
	"strings"
	"fmt"
)

type GraphiteClient struct {
	Api string
	*graphite.Graphite
}

type ApiResponse []ApiMetric

type ApiMetric struct {
	Target	string `json: "target"`
	Datapoints []Datapoint `json: "datapoints"`
}

type Datapoint [2]interface{}

func (d *Datapoint) Timestamp() int64 {
	return d[1].(int64)
}

func (d *Datapoint) Value() string {
	if d[0] == nil {
		return ""
	}
	f := d[0].(float64)
	s := strconv.FormatFloat(f, 'f', -1, 64)

	return s
}

func NewGraphiteClient(host string, port int, apiHost string, protocol string, prefix string) (*GraphiteClient, error) {
	var client *GraphiteClient
	g, err :=  graphite.GraphiteFactory(protocol, host, port, prefix)
	if err != nil {
		return nil, err
	}
	client = &GraphiteClient{Api: apiHost, Graphite: g}
	return client, nil
}

func (client *GraphiteClient) SendMetricToGraphite(name string, value string) error{

	log.Printf("Loaded Graphite connection: %#v", client)
	log.Printf("Original graphite client: %#v", client.Graphite)

	err := client.SimpleSend(name, value)
	return err
}

func (client *GraphiteClient) GetMetricFromGraphite(name string, from int) (ApiResponse, error) {
	httpClient := &http.Client{}

	apiEndpoint := strings.Join([]string{"http://", client.Api, "/render"}, "")
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		log.Printf("Error while creating a new GET request for endpoint %s", apiEndpoint)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("target", strings.Join([]string{client.Prefix, name}, "."))
	q.Add("format", "json")
	q.Add("from", fmt.Sprintf("-%ds", from))
	req.URL.RawQuery = q.Encode()
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("Error while attempting a request to %s endpoint", apiEndpoint)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading the response body for request to %s endpoint", apiEndpoint)
		return nil, err
	}

	decResponse, err := DecodeApiResponse(body)
	if err != nil {
		log.Printf("Error while decoding the json api response for request to %s", apiEndpoint)
		return nil, err
	}

	return decResponse, nil
}

func DecodeApiResponse(respBody []byte) (ApiResponse, error) {
	var resp ApiResponse
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
