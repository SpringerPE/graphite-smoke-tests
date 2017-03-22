package graphite

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/SpringerPE/graphite-smoke-tests/smoke"
	"time"
	"fmt"
	"strconv"
	"log"
)


var _ = Describe("Graphite:", func() {

	testConfig := smoke.GetConfig()
	testPrefix := "smoke_tests"

	var (
		gclient *smoke.GraphiteClient
	)

	if testConfig.TcpEnabled {
		Context("send single tcp metric value", func() {

	    		var err error
	    		protocol := "tcp"

			gclient, err = smoke.NewGraphiteClient(
				testConfig.Host, 
				testConfig.Port,
				testConfig.Api,
				protocol,
				testPrefix,
			)

			if err != nil {
				Fail(err.Error())
			}

			It("can be sent over tcp and retrieved", func() {
				runSendAndRetrieveTest(gclient, protocol)
			})
		})
	}

	if testConfig.UdpEnabled {
		Context("send single udp metric value", func() {


	    		var err error
	    		protocol := "udp"

			gclient, err = smoke.NewGraphiteClient(
				testConfig.Host, 
				testConfig.Port,
				testConfig.Api,
				protocol,
				testPrefix,
			)

			if err != nil {
				Fail(err.Error())
			}

			It("can be sent over udp and retrieved", func() {
		    		if ! testConfig.UdpEnabled {
		        			Skip("udp not enabled, skip test")
		    		}
				runSendAndRetrieveTest(gclient, protocol)
			})
		})
	}
})

func runSendAndRetrieveTest(gclient *smoke.GraphiteClient, metricBase string) {

	var (
		apiResponse smoke.ApiResponse
		err error
		elapsed int
	)

	maxAttempts := 60
	retryDelay := 10000
	from := retryDelay / 1000
	found := false
	value := strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10)
	metricName := fmt.Sprintf("%s.%s", metricBase , value)
	expectedOutput := fmt.Sprintf("datapoint: %s, for metric %s", value, metricName)

	//Start the test
	err = gclient.SendMetricToGraphite(metricName, value)
	if err != nil {
		Fail(fmt.Sprintf("Error while executing the test: %s\n", err.Error()))
	}
	elapsed = 0
	Loop:
		for j := 0; j < maxAttempts; j++ {
			apiResponse, err = gclient.GetMetricFromGraphite(metricName, from)
			log.Printf("Response from the api server: %v\n", apiResponse)
			log.Printf("Expected output: %s\n", expectedOutput)
			if err != nil {
				break Loop
			}
			for _, metric := range(apiResponse) {
				for _, datapoint := range(metric.Datapoints) {
					fmt.Printf("%v\n", datapoint)
					if datapoint.Value() == value {
						found = true
						break Loop
					}
				}
				log.Printf("%s not found in: %v\n", expectedOutput, apiResponse)
			}
			time.Sleep(time.Duration(retryDelay) * time.Millisecond)
			elapsed += retryDelay
			log.Printf("Attempt #%d, sec elapsed: %d\n", j+1, elapsed/1000)
		}

	Expect(found).To(BeTrue(), fmt.Sprintf("Wanted to see '%s' (datapoint found) in %d attempts, but didn't\n", expectedOutput, maxAttempts))
}
