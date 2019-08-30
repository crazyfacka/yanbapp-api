package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestYanbappApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "YanbappApi Suite")
}
