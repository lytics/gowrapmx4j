package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func (m MX4J) QueryMX4JBean(mm MX4JMetric) (*Bean, error) {
	query := fmt.Sprintf("mbean?objectname=%s&template=identity", mm.ObjectName)
	fullQuery := m.hostAddr + query
	log.Debug(fullQuery)

	httpResp, err := http.Get(fullQuery)
	if err != nil {
		log.Errorf("Failed to get response from mx4j: %#v", err)
		return nil, err
	}
	mb, err := getBeans(httpResp.Body, beanUnmarshal)
	if err != nil {
		log.Errorf("Error getting attribute: %s %s %s", mm.ObjectName, mm.Format, mm.Attribute)
		return nil, err
	}
	return mb, err
}

//Handles reading of the http.Body and passes bytes of io.ReadCloser
//to getAttrUnmarshal() for unmarshaling XML.
func getBeans(httpBody io.ReadCloser, unmarshalFunc func([]byte) (*Bean, error)) (*Bean, error) {
	xmlBytes, err := ioutil.ReadAll(httpBody)
	if err != nil {
		log.Errorf("Failed to read http response: %#v", err)
		return nil, err
	}

	return unmarshalFunc(xmlBytes)
}

//Unmarshals XML and returns an Bean struct
func beanUnmarshal(xmlBytes []byte) (*Bean, error) {
	var mb Bean
	err := xml.Unmarshal([]byte(xmlBytes), &mb)
	if err != nil {
		log.Errorf("Failed to Unmarshal xml: %#v", err)
		log.Errorf("Bytes failed to be unmarshalled: \n%s", xmlBytes)
		return nil, err
	}
	return &mb, nil
}
