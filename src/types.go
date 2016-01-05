package gowrapmx4j

type MX4JMetric struct {
	HumanName  string
	ObjectName string
	Format     string
	Attribute  string
	ValFunc    func(*MBean) map[string]string
	MBean      *MBean
}

func NewMX4JMetric(hname, objname, format, attr string) MX4JMetric {
	return MX4JMetric{HumanName: hname, ObjectName: objname, Format: format, Attribute: attr}
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
