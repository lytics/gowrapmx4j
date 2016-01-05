package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

//Struct representing MX4J address to query
type MX4J struct {
	Host     string
	Port     string
	hostAddr string
}

func (m *MX4J) Init() {
	m.hostAddr = fmt.Sprintf("http://%s:%s/", m.Host, m.Port)
}

// Queries MX4J to get an attribute's data, returns MBean struct or error
// equivalent to http://hostname:port/getattribute?queryargs...
// eg: "http://localhost:8081/getattribute?objectname=org.apache.cassandra.metrics:type=ColumnFamily,keyspace=yourkeyspace,scope=node,name=ReadLatency&format=array&attribute=Max&template=identity"
func (m MX4J) QueryGetAttributes(objectname, format, attribute string) (*MBean, error) {

	query := fmt.Sprintf("getattribute?objectname=%s&format=%s&attribute=%s&template=identity", objectname, format, attribute) //template?
	fullQuery := m.hostAddr + query
	log.Debug(fullQuery)

	httpResp, err := http.Get(fullQuery)
	if err != nil {
		log.Errorf("Failed to get response from mx4j: %#v", err)
		return nil, err
	}
	return getAttributes(httpResp.Body, getAttrUnmarshal)
}

func (m MX4J) QueryMX4JMetric(mm MX4JMetric) (*MBean, error) {
	query := fmt.Sprintf("getattribute?objectname=%s&format=%s&attribute=%s&template=identity", mm.ObjectName, mm.Format, mm.Attribute) //template?
	fullQuery := m.hostAddr + query
	log.Debug(fullQuery)

	httpResp, err := http.Get(fullQuery)
	if err != nil {
		log.Errorf("Failed to get response from mx4j: %#v", err)
		return nil, err
	}
	return getAttributes(httpResp.Body, getAttrUnmarshal)
}

//Handles reading of the http.Body and passes bytes of io.ReadCloser
//to getAttrUnmarshal() for unmarshaling XML.
func getAttributes(httpBody io.ReadCloser, unmarshalFunc func([]byte) (*MBean, error)) (*MBean, error) {
	xmlBytes, err := ioutil.ReadAll(httpBody)
	if err != nil {
		log.Errorf("Failed to read http response: %#v", err)
		return nil, err
	}

	mb, err := unmarshalFunc(xmlBytes)

	return mb, nil
}

//Unmarshals XML and returns an MBean struct
func getAttrUnmarshal(xmlBytes []byte) (*MBean, error) {
	var mb MBean
	err := xml.Unmarshal([]byte(xmlBytes), &mb)
	if err != nil {
		log.Errorf("Failed to Unmarshal xml: %#v", err)
		log.Errorf("Bytes failed to be unmarshalled: \n%s", xmlBytes)
		return nil, err
	}
	return &mb, nil
}
