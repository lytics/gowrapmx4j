package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type MX4JData interface {
	QueryMX4J(m MX4J, mm MX4JMetric) (MX4JData, error)
}

type MX4JMetric struct {
	HumanName  string
	ObjectName string
	Format     string
	Attribute  string
	ValFunc    func(MX4JData) map[string]string
	MetricFunc func(MX4JData, string)
	Data       MX4JData
}

func NewMX4JMetric(hname, objname, format, attr string) MX4JMetric {
	return MX4JMetric{HumanName: hname, ObjectName: objname, Format: format, Attribute: attr}
}

// AKA Whole Bean of data
type Bean struct {
	XMLName    xml.Name        `xml:"MBean"`
	ObjectName string          `xml:"objectname,attr"`
	ClassName  string          `xml:"classname,attr"`
	Attributes []MX4JAttribute `xml:"Attribute"`
}

func (b Bean) AttributeMap() map[string]MX4JAttribute {
	attrMap := make(map[string]MX4JAttribute)
	for _, a := range b.Attributes {
		attrMap[a.Name] = a
	}
	return attrMap
}

func (b Bean) QueryMX4J(m MX4J, mm MX4JMetric) (MX4JData, error) {
	query := fmt.Sprintf("mbean?objectname=%s&template=identity", mm.ObjectName)
	fullQuery := m.hostAddr + query
	log.Debug(fullQuery)

	httpResp, err := http.Get(fullQuery)
	defer httpResp.Body.Close()
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

/*Example XML
<?xml version="1.0" encoding="UTF-8"?>
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=yourkeyspace,scope=node,name=ReadLatency">
  <Attribute classname="double" isnull="false" name="Max" value="0.0"/>
</MBean>
*/

type MBean struct {
	ObjectName string        `xml:"objectname,attr"`
	ClassName  string        `xml:"classname,attr"`
	Attribute  MX4JAttribute `xml:"Attribute"`
}

func (mbean MBean) QueryMX4J(m MX4J, mm MX4JMetric) (MX4JData, error) {
	query := fmt.Sprintf("getattribute?objectname=%s&format=%s&attribute=%s&template=identity", mm.ObjectName, mm.Format, mm.Attribute) //template?
	fullQuery := m.hostAddr + query
	log.Debug(fullQuery)

	httpResp, err := http.Get(fullQuery)
	defer httpResp.Body.Close()
	if err != nil {
		log.Errorf("Failed to get response from mx4j: %#v", err)
		return nil, err
	}
	mb, err := getAttributes(httpResp.Body, getAttrUnmarshal)
	if err != nil {
		log.Errorf("Error getting attribute: %s %s %s", mm.ObjectName, mm.Format, mm.Attribute)
		return nil, err
	}
	return *mb, err
}

type MX4JAttribute struct {
	Classname string  `xml:"classname,attr"`
	Name      string  `xml:"name,attr"`
	Value     string  `xml:"value,attr"`
	Map       MX4JMap `xml:"Map"`
}

type MX4JMap struct {
	Length   string        `xml:"length,attr"`
	Elements []MX4JElement `xml:"Element"`
}

type MX4JElement struct {
	Key     string `xml:"key,attr"`
	Element string `xml:"element,attr"` //Known as 'Value' to the rest of the world
	Index   string `xml:"index,attr"`
}
