package gowrapmx4j

import "sync"

var registry = make(map[string]*MX4JMetric)
var reglock = &sync.RWMutex{}

func RegistrySet(mm MX4JMetric, mb *MBean) {
	reglock.Lock()
	defer reglock.Unlock()

	mm.MBean = mb
	registry[mm.HumanName] = &mm
}

func RegistryGet(humanName string) *MX4JMetric {
	reglock.RLock()
	defer reglock.RUnlock()

	return registry[humanName]
}

func RegistryBeans() map[string]*MBean {
	reglock.RLock()
	defer reglock.RUnlock()

	beans := make(map[string]*MBean)
	for hname, mm := range registry {
		beans[hname] = mm.MBean
	}
	return beans
}

func RegistryGetAll() []MX4JMetric {
	reglock.RLock()
	defer reglock.RUnlock()
	metrics := make([]MX4JMetric, 0, 0)
	for _, mm := range registry {
		metrics = append(metrics, *mm)
	}
	return metrics
}
