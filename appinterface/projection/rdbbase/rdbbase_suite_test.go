package rdbbase_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRDbBase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RDbBase Suite")
}
