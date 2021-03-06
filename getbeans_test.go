package gowrapmx4j

import (
	"encoding/xml"
	"fmt"
	"testing"
)

var rawBean = `<?xml version="1.0" encoding="UTF-8"?>
<MBean classname="com.yammer.metrics.reporting.JmxReporter$Timer" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.metrics:type=ColumnFamily,keyspace=yourkeyspace,scope=node,name=RangeLatency">
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="50thPercentile" strinit="true" type="double" value="2.4563736922E7"/>
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="75thPercentile" strinit="true" type="double" value="4.9114161986E7"/>
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="95thPercentile" strinit="true" type="double" value="4.9114161986E7"/>
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="98thPercentile" strinit="true" type="double" value="4.9114161986E7"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="999thPercentile" strinit="true" type="double" value="4.9114161986E7"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="99thPercentile" strinit="true" type="double" value="4.9114161986E7"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="Count" strinit="true" type="long" value="23"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="EventType" strinit="true" type="java.lang.String" value="calls"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="FifteenMinuteRate" strinit="true" type="double" value="4.44659081257E-313"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="FiveMinuteRate" strinit="true" type="double" value="1.4821969375E-313"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="LatencyUnit" strinit="false" type="java.util.concurrent.TimeUnit" value="MICROSECONDS"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="Max" strinit="true" type="double" value="4.9114161986E7"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="Mean" strinit="true" type="double" value="4411217.701"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="MeanRate" strinit="true" type="double" value="2.3773582427257664E-6"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="Min" strinit="true" type="double" value="13311.858"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="OneMinuteRate" strinit="true" type="double" value="2.964393875E-314"/><Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="RateUnit" strinit="false" type="java.util.concurrent.TimeUnit" value="SECONDS"/>
  <Attribute availability="RO" description="Attribute exposed for management" isnull="false" name="StdDev" strinit="true" type="double" value="1.0582224553531626E7"/>
  <Operation description="Operation exposed for management" impact="unknown" name="objectName" return="javax.management.ObjectName"/>
  <Operation description="Operation exposed for management" impact="unknown" name="values" return="[D"/>
</MBean>`

func TestUnMarshalXML(t *testing.T) {

	x := Bean{ObjectName: "neh"}

	err := xml.Unmarshal([]byte(rawBean), &x)
	if err != nil {
		t.Errorf("Error unmarshalling: %v\n", err)
	}

	//fmt.Printf("%#v\n", x)
	if len(x.Attributes) < 10 {
		t.Errorf("Error unmarshalling attributes")
	}

}

func TestUnmarshallingFunctions(t *testing.T) {
	b, err := beanUnmarshal([]byte(rawBean))
	if err != nil {
		t.Errorf("Error unmarshalling: %v\n", err)
	}

	//fmt.Printf("%#v\n", b)
	for _, x := range b.Attributes {
		fmt.Printf("%s %s\n", x.Name, x.Value)
	}
	if len(b.Attributes) < 5 {
		t.Errorf("Seems like number of Attributes unmarshalled is a little low..\n%#v\n", b.Attributes)
	}
}

func TestInterfaceAttributes(t *testing.T) {
	input := `<MBean classname="org.apache.cassandra.gms.FailureDetector" description="Information on the management interface of the MBean" objectname="org.apache.cassandra.net:type=FailureDetector">
	<Attribute classname="java.util.Map" isnull="false" name="SimpleStates">
		<Map length="1">
			<Element element="UP" elementclass="java.lang.String" index="0" key="/127.0.0.1" keyclass="java.lang.String"/>
		</Map>
	</Attribute>
</MBean>`

	var bean Bean
	err := xml.Unmarshal([]byte(input), &bean)
	if err != nil {
		t.Errorf("Error unmarshalling MX4J data: %v", err)
	}
	_, err = beanUnmarshal([]byte(rawBean))
	if err != nil {
		t.Errorf("Error unmarshalling: %v\n", err)
	}

	// Enable for helpfull debugging
	//katoptron.Display("MBean", bean)
	//katoptron.Display("MBean", *b)
}
