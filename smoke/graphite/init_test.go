package graphite

import (
	"testing"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo"
)


func TestSmokeTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Graphite Suite")

}