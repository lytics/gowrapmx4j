package gowrapmx4j

import (
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestContextGet(t *testing.T) {
	to := time.Second * 1
	ctx, _ := context.WithTimeout(context.Background(), to)

	mm := NewMX4JMetric("NodeStatus", "org.apache.cassandra.net:type=FailureDetector", "map", "SimpleStates")
	RegistrySet(mm, nil)

	m := HTTPRegistryGetAll(ctx)

	if m == nil {
		t.Errorf("error retrieving metrics from registry")
	}
	if m != nil && len(*m) != 1 {
		t.Errorf("Registry metric set was not returned")
	}
}

func TestContextGetFailure(t *testing.T) {
	to := time.Second * 0
	ctx, _ := context.WithTimeout(context.Background(), to)

	mm := NewMX4JMetric("NodeStatus", "org.apache.cassandra.net:type=FailureDetector", "map", "SimpleStates")
	RegistrySet(mm, nil)

	m := HTTPRegistryGetAll(ctx)

	if m != nil {
		t.Errorf("metrics should be nil")
	}
}
