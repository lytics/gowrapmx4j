package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type MX4J struct {
	Host     string
	Port     string
	hostAddr string
}

func (m *MX4J) Init() {
	m.hostAddr = fmt.Sprintf("http://%s:%s/", m.Host, m.Port)
}

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

	var mb MBean
	err = xml.Unmarshal([]byte(xmlBytes), &mb)
	if err != nil {
		log.Errorf("Failed to Unmarshal xml: %#v", err)
		return nil, err
	}

	return &mb, nil
}
