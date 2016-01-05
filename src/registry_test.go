package gowrapmx4j

import "testing"

func TestMetricSet(t *testing.T) {

	//func NewMX4JMetric(hname, objname, format, attr string) MX4JMetric {
	mm := NewMX4JMetric("NodeStatus", "org.apache.cassandra.net:type=FailureDetector", "map", "SimpleStates")
	RegistrySet(mm, nil)

	if len(registry) != 1 {
		t.Errorf("Registry was not updated")
	}

	metrics := RegistryGetAll()
	if len(metrics) != 1 {
		t.Errorf("Not all metrics returned")
	}
}

func TestMultiMetricSet(t *testing.T) {

	//func NewMX4JMetric(hname, objname, format, attr string) MX4JMetric {
	mm := NewMX4JMetric("NodeStatus", "org.apache.cassandra.net:type=FailureDetector", "map", "SimpleStates")
	RegistrySet(mm, nil)
	RegistrySet(mm, nil)
	RegistrySet(mm, nil)
	RegistrySet(mm, nil)

	if len(registry) != 1 {
		t.Errorf("Registry was not updated")
	}

	metrics := RegistryGetAll()
	if len(metrics) != 1 {
		t.Errorf("Not all metrics returned")
	}
	RegistrySet(mm, nil)
	RegistrySet(mm, nil)

	metrics = RegistryGetAll()
	if len(metrics) != 1 {
		t.Errorf("Not all metrics returned")
	}
}

func TestMultiMetricGet(t *testing.T) {

	//func NewMX4JMetric(hname, objname, format, attr string) MX4JMetric {
	mm := NewMX4JMetric("NodeStatus", "org.apache.cassandra.net:type=FailureDetector", "map", "SimpleStates")
	RegistrySet(mm, nil)

	hname := "NodeStatus"
	m := RegistryGet(hname)
	if m == nil {
		t.Errorf("%s not found in registry", hname)
	}
	//fmt.Printf("%#v\n", m)
	if m.HumanName != hname {
		t.Errorf("Wrong retristry metric returned")
	}
}
