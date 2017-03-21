package smoke

import (
	"io/ioutil"
	"encoding/json"
	graphite "github.com/marpaia/graphite-golang"
	"net/http"
	"log"
	"strconv"
	"strings"
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

func NewGraphiteClient(host string, port int, apiHost string, protocol string, prefix string) *GraphiteClient {
	var client *GraphiteClient

	g, _ :=  graphite.GraphiteFactory(protocol, host, port, prefix)
	client = &GraphiteClient{Api: apiHost, Graphite: g}
	return client
}

func (client *GraphiteClient) SendMetricToGraphite(name string, value string) int {

	log.Printf("Loaded Graphite connection: %#v", client)
	client.SimpleSend(name, value)
	return 0
}

func (client *GraphiteClient) GetMetricFromGraphite(name string) (ApiResponse, error) {
	httpClient := &http.Client{}

	api_endpoint := strings.Join([]string{"http://", client.Api, "/render"}, "")
	req, err := http.NewRequest("GET", api_endpoint, nil)
	if err != nil {
		log.Printf("Error while creating a new GET request for endpoint %s", api_endpoint)
		return nil, err
	}
	q := req.URL.Query()
	q.Add("target", strings.Join([]string{client.Prefix, name}, "."))
	q.Add("format", "json")
	q.Add("from", "-5min")
	req.URL.RawQuery = q.Encode()
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error while attempting a request to %s endpoint", api_endpoint)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading the response body for request to %s endpoint", api_endpoint)
		return nil, err
	}
	dec_response, err := DecodeApiResponse(body)
	if err != nil {
		log.Printf("Error while decoding the json api response for request to %s", api_endpoint)
		return nil, err
	}

	return dec_response, nil
}

func DecodeApiResponse(respBody []byte) (ApiResponse, error) {
	var resp ApiResponse
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
