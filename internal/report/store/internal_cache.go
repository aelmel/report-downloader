package store

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type ReportStore interface {
	AddReport(id string)
	GetReportsId() []string
	RemoveReport(id string)
}

type store struct {
	cache  map[string]int
	logger *logrus.Logger
	sync.Mutex
}

const (
	newReport = iota
	inProgress
)

func NewReportStore(logger *logrus.Logger) ReportStore {
	cache := make(map[string]int)

	return &store{
		cache:  cache,
		logger: logger,
		Mutex:  sync.Mutex{},
	}
}

func (s *store) AddReport(id string) {
	s.Lock()
	defer s.Unlock()
	s.cache[id] = newReport

}

func (s *store) GetReportsId() []string {
	s.Lock()
	defer s.Unlock()

	keys := make([]string, len(s.cache))
	for k, v := range s.cache {
		if v == newReport {
			keys = append(keys, k)
			s.cache[k] = inProgress
		}
	}

	return keys
}

func (s *store) RemoveReport(id string) {
	s.Lock()
	defer s.Unlock()

	delete(s.cache, id)
}
