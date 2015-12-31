package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
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
// eg: "http://localhost:8081/getattribute?objectname=org.apache.cassandra.metrics:type=ColumnFamily,keyspace=lio4,scope=node,name=ReadLatency&format=array&attribute=Max&template=identity"
func (m MX4J) GetAttributes(objectname, format, attribute string) (*MBean, error) {
	query := fmt.Sprintf("getattribute?objectname=%s&format=array&attribute=Max&template=identity", objectname, format, attribute) //template?

	httpResp, err := http.Get(m.hostAddr + query)
	if err != nil {
		log.Errorf("Failed to get response from mx4j: %#v", err)
		return nil, err
	}

	xmlBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Errorf("Failed to read http response: %#v", err)
		return nil, err
	}

	mb, err := GetAttrUnmarshal(xmlBytes)

	return mb, nil
}

func GetAttrUnmarshal(xmlBytes []byte) (*MBean, error) {
	var mb MBean
	err := xml.Unmarshal([]byte(xmlBytes), &mb)
	if err != nil {
		log.Errorf("Failed to Unmarshal xml: %#v", err)
		return nil, err
	}
	return &mb, nil
}
