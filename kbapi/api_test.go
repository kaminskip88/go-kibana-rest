package kbapi

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type KBAPITestSuite struct {
	suite.Suite
	client *resty.Client
	*API
}

func (s *KBAPITestSuite) SetupSuite() {

	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	address := os.Getenv("KIBANA_URL")
	username := os.Getenv("KIBANA_USERNAME")
	password := os.Getenv("KIBANA_PASSWORD")

	if address == "" {
		panic("You need to put kibana url on environment variable KIBANA_URL. If you need auth, you can use KIBANA_USERNAME and KIBANA_PASSWORD")
	}

	restyClient := resty.New().
		SetHostURL(address).
		SetBasicAuth(username, password).
		SetHeader("kbn-xsrf", "true").
		SetHeader("Content-Type", "application/json")

	s.client = restyClient
	s.API = New(restyClient)

	// Wait kb is online
	isOnline := false
	nbTry := 0
	for isOnline == false {
		_, err := s.API.KibanaSpaces.List()
		if err == nil {
			isOnline = true
		} else {
			time.Sleep(5 * time.Second)
			if nbTry == 10 {
				panic(fmt.Sprintf("We wait 50s that Kibana start: %s", err))
			}
			nbTry++
		}
	}

	// Create kibana space
	space := &KibanaSpace{
		ID:   "testacc",
		Name: "testacc",
	}
	_, err := s.API.KibanaSpaces.Create(space)
	if err != nil {
		if err.(APIError).Code != 409 {
			panic(err)
		}
	}

	// Generate test data

	// currentTime := time.Now()

	// b := `{"mappings": {"properties": {"@timestamp": {"type": "date" }}}}`
	// s.client.R().SetBody(b).
	// 	Put(fmt.Sprintf("http://localhost:9200/logs-%s/", currentTime.Format("2006.01.02")))

	// for i := 1; i <= 100; i++ {
	// 	date := currentTime.Add(time.Second * time.Duration(10*i))
	// 	d := testDoc{
	// 		Timestamp: date.UnixMilli(),
	// 		Message:   fmt.Sprintf("Test message index %d", i),
	// 	}
	// 	s.client.R().SetBody(d).
	// 		Post(fmt.Sprintf("http://localhost:9200/logs-%s/_doc", currentTime.Format("2006.01.02")))
	// }

}

type testDoc struct {
	Timestamp int64  `json:"@timestamp"`
	Message   string `json:"message"`
}

func (s *KBAPITestSuite) SetupTest() {

	// Do somethink before each test

}

func TestKBAPITestSuite(t *testing.T) {
	suite.Run(t, new(KBAPITestSuite))
}
