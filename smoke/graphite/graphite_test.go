package graphite

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/SpringerPE/graphite-smoke-tests/smoke"
	"time"
	"fmt"
	"strconv"
)


var _ = Describe("Graphite:", func() {

	testConfig := smoke.GetConfig()
	testPrefix := "smoke_tests"

	var (
		gclient *smoke.GraphiteClient
	)


	BeforeEach(func() {


		})

	AfterEach(func() {


		})

	Context("send single tcp metric value", func() {

		gclient = smoke.NewGraphiteClient(
			testConfig.Host, 
			testConfig.Port,
			testConfig.Api,
			testConfig.Proto,
			testPrefix,
		)

		It("can be sent over tcp and retrieved", func() {
			runSendAndRetrieveTest(gclient, "tcp_metric")
		})
	})
})

func sendMetric() int{
	return 0
}

func runSendAndRetrieveTest(gclient *smoke.GraphiteClient, metric_name string) {

	var api_response smoke.ApiResponse
	var err error

	maxAttempts := 60
	found := false

	value := strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10)
	gclient.SendMetricToGraphite(metric_name, value)

	expectedOutput := fmt.Sprintf("datapoint: %s, for metric %s", value, metric_name)

	Loop:
		for j := 0; j < maxAttempts; j++ {
			api_response, err = gclient.GetMetricFromGraphite(metric_name)
			if err != nil {
				break Loop
			}
			for _, metric := range(api_response) {
				for _, datapoint := range(metric.Datapoints) {
					if datapoint.Value() == value {
						found = true
						break Loop
					}
				}
			}
			time.Sleep(5000 * time.Millisecond)
		}

	Expect(found).To(BeTrue(), fmt.Sprintf("Wanted to see '%s' (datapoint found) in %d attempts, but didn't", expectedOutput, maxAttempts))
}
