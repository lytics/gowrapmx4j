package gowrapmx4j

import "encoding/xml"

/*Example XML
<?xml version="1.0" encoding="UTF-8"?>
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=lio4,scope=node,name=ReadLatency">
  <Attribute classname="double" isnull="false" name="Max" value="0.0"/>
</MBean>
*/

type MBean struct {
	XMLName    xml.Name      `xml:"MBean"`
	ObjectName string        `xml:"objectname,attr"`
	ClassName  string        `xml:"classname,attr"`
	Attribute  MX4JAttribute `xml:"Attribute"`
}

type MX4JAttribute struct {
	XMLName   xml.Name `xml:"Attribute"`
	Classname string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Value     string   `xml:"value,attr"`
}
