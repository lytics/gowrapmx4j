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
