package smoke

import (
	"io/ioutil"
	"encoding/json"
	graphite "github.com/marpaia/graphite-golang"
	"net/http"
	"log"
	"fmt"
	//"time"
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

type Datapoint [2]float64

func (d *Datapoint) Timestamp() int64 {
	return int64(d[1])
}

func (d *Datapoint) Value() float64 {
	return d[0]
}

func NewGraphiteClient(host string, port int, apiHost string) *GraphiteClient {
	var client *GraphiteClient

	g, _ :=  graphite.NewGraphite(host, port)

	client = &GraphiteClient {Api: apiHost, Graphite: g}
	return client
}

func (client *GraphiteClient) SendMetricToGraphite() int {

	log.Printf("Loaded Graphite connection: %#v", client)
	client.SimpleSend("test.smoke.graphite", "1")
	return 0
}

func (client *GraphiteClient) GetMetricFromGraphite() []ApiMetric {
	httpClient := &http.Client{}

	req, _ := http.NewRequest("GET", "http://" + client.Api + "/render", nil)
	q := req.URL.Query()
	q.Add("target", "test.smoke.graphite")
	q.Add("format", "json")
	q.Add("from", "-5min")
	req.URL.RawQuery = q.Encode()
	resp, _ := httpClient.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return DecodeApiResponse(body)
}

func DecodeApiResponse(respBody []byte) []ApiMetric{
	var resp ApiResponse
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		
	}
	return resp
}
