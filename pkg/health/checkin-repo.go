package health

import (
	"errors"
	"sync"

	"github.com/heptiolabs/healthcheck"
)

type HealthRecord struct {
	Name            string
	Healthy         bool
	UnhealthyReason string
}

var (
	repo sync.Map
)

func CheckIn(record HealthRecord) {
	// Store an item in the map.
	repo.Store(record.Name, record)
}

func GetHealthRecord(name string) (*HealthRecord, bool) {
	result, ok := repo.Load(name)
	if ok {
		r := result.(HealthRecord)
		return &r, ok
	}
	return nil, false
}

func CreateHealthCheck(name string) healthcheck.Check {
	var check healthcheck.Check = func() error {
		record, ok := GetHealthRecord(name)
		if ok && record != nil {
			if record.Healthy {
				return nil
			}
			return errors.New(record.UnhealthyReason)
		}
		return errors.New("No entry for this healthcheck was found")
	}
	return check
}
