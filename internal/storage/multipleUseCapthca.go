package storage

import (
	"container/list"
	"github.com/mojocn/base64Captcha"
	"strings"
	"sync"
	"time"
)

type captchaWithTimes struct {
	value string
	times int
}

type idByTimeValue struct {
	timestamp time.Time
	id        string
}

// multipleUseCaptcha is a multiple use version of memory store for captcha ids and their values.
type multipleUseCaptcha struct {
	sync.RWMutex
	digitsById map[string]captchaWithTimes
	idByTime   *list.List
	// Number of items stored since last collection.
	numStored int
	// Number of saved items that triggers collection.
	collectNum int
	// Expiration time of captchas.
	expiration time.Duration
	useTimes   int
}

// NewMultipleUseCaptcha returns a new multiple-use memory store for captchas with the
// given collection threshold and expiration time (duration). The returned
// store must be registered with SetCustomStore to replace the default one.
func NewMultipleUseCaptcha(collectNum int, expiration time.Duration, useTimes int) base64Captcha.Store {
	s := new(multipleUseCaptcha)
	s.digitsById = make(map[string]captchaWithTimes)
	s.idByTime = list.New()
	s.collectNum = collectNum
	s.expiration = expiration
	s.useTimes = useTimes
	return s
}

func (s *multipleUseCaptcha) Set(id string, value string) error {
	s.Lock()
	s.digitsById[id] = captchaWithTimes{
		value: value,
		times: s.useTimes,
	}
	s.idByTime.PushBack(idByTimeValue{time.Now(), id})
	s.numStored++
	s.Unlock()
	if s.numStored > s.collectNum {
		go s.collect()
	}
	return nil
}

func (s *multipleUseCaptcha) Verify(id, answer string, clear bool) bool {
	if id == "" || answer == "" {
		return false
	}
	v := s.Get(id, clear)
	return strings.EqualFold(v, answer)
}

func (s *multipleUseCaptcha) Get(id string, clear bool) string {
	if !clear {
		// When we don't need to clear captcha, acquire read lock.
		s.RLock()
		defer s.RUnlock()
	} else {
		s.Lock()
		defer s.Unlock()
	}
	cwt, ok := s.digitsById[id]
	if !ok {
		return ""
	}
	if clear {
		cwt.times--
		if cwt.times <= 0 {
			delete(s.digitsById, id)
			return ""
		}
		s.digitsById[id] = cwt
	}
	return cwt.value
}

func (s *multipleUseCaptcha) collect() {
	now := time.Now()
	s.Lock()
	defer s.Unlock()
	for e := s.idByTime.Front(); e != nil; {
		e = s.collectOne(e, now)
	}
}

func (s *multipleUseCaptcha) collectOne(e *list.Element, specifyTime time.Time) *list.Element {

	ev, ok := e.Value.(idByTimeValue)
	if !ok {
		return nil
	}

	if ev.timestamp.Add(s.expiration).Before(specifyTime) {
		delete(s.digitsById, ev.id)
		next := e.Next()
		s.idByTime.Remove(e)
		s.numStored--
		return next
	}
	return nil
}
