package storage

import "CloudPhoto/config"

type MultiUseStore struct {
	data map[string]struct {
		value string
		count int
	}
}

func (s *MultiUseStore) Set(id, value string) error {
	s.data[id] = struct {
		value string
		count int
	}{value: value, count: config.Get().App.CaptchaUseTimes}
	return nil
}

func (s *MultiUseStore) Get(id string, clear bool) (value string) {
	item, ok := s.data[id]
	if !ok {
		return ""
	}
	if clear {
		delete(s.data, id)
	}
	return item.value
}

func (s *MultiUseStore) Verify(id, answer string, clear bool) bool {
	item, ok := s.data[id]
	if !ok {
		return false
	}
	if item.value == answer {
		if clear {
			item.count = 0
		} else {
			item.count--
		}
		if item.count <= 0 {
			delete(s.data, id)
		} else {
			s.data[id] = item
		}
		return true
	}
	return false
}
