package main

import (
	mt "MediaUnlockTest/pkg/providers"
	core "MediaUnlockTest/pkg/core"
	"net/http"
	"sync"
	"time"
)

var (
	MUL  bool
	HK   bool
	TW   bool
	JP   bool
	KR   bool
	NA   bool
	SA   bool
	EU   bool
	AFR  bool
	SEA  bool
	OCEA bool
	AI   bool
	Conc uint64 = 0
	sem  chan struct{}
)

type TEST struct {
	Client  http.Client
	Results []*result
	Wg      *sync.WaitGroup
}

func NewTest() *TEST {
	t := &TEST{
		Client:  core.AutoHttpClient,
		Results: make([]*result, 0),
		Wg:      &sync.WaitGroup{},
	}
	if Conc > 0 {
		sem = make(chan struct{}, Conc) // 初始化带缓冲的通道
	}
	return t
}

func (T *TEST) Check() bool {
	if MUL {
		T.Multination()
	}
	if HK {
		T.HongKong()
	}
	if TW {
		T.Taiwan()
	}
	if JP {
		T.Japan()
	}
	if KR {
		T.Korea()
	}
	if NA {
		T.NorthAmerica()
	}
	if SA {
		T.SouthAmerica()
	}
	if EU {
		T.Europe()
	}
	if AFR {
		T.Africa()
	}
	if SEA {
		T.SouthEastAsia()
	}
	if OCEA {
		T.Oceania()
	}
	if AI {
		T.AI()
	}

	ch := make(chan struct{})
	go func() {
		defer close(ch)
		T.Wg.Wait()
	}()
	select {
	case <-ch:
		return false
	case <-time.After(30 * time.Second):
		return true
	}
}

type result struct {
	Type  string
	Name  string
	Value core.Result
}

func (T *TEST) excute(Name string, F func(client http.Client) core.Result) {
	r := &result{Name: Name}
	T.Results = append(T.Results, r)
	T.Wg.Add(1)
	go func() {
		if Conc > 0 {
			sem <- struct{}{}
			defer func() {
				<-sem
				T.Wg.Done()
			}()
		} else {
			defer T.Wg.Done()
		}
		r.Value = F(T.Client)
	}()
}

func (T *TEST) executeList(tests []mt.TestItem) {
	for _, test := range tests {
		if test.Func != nil {
			T.excute(test.Name, test.Func)
		}
	}
}

func (T *TEST) Multination() {
	T.executeList(mt.GlobeTests)
}

func (T *TEST) HongKong() {
	T.executeList(mt.HongKongTests)
}

func (T *TEST) Taiwan() {
	T.executeList(mt.TaiwanTests)
}

func (T *TEST) Japan() {
	T.executeList(mt.JapanTests)
}

func (T *TEST) Korea() {
	T.executeList(mt.KoreaTests)
}

func (T *TEST) NorthAmerica() {
	T.executeList(mt.NorthAmericaTests)
}

func (T *TEST) SouthAmerica() {
	T.executeList(mt.SouthAmericaTests)
}

func (T *TEST) Europe() {
	T.executeList(mt.EuropeTests)
}

func (T *TEST) Africa() {
	T.executeList(mt.AfricaTests)
}

func (T *TEST) SouthEastAsia() {
	T.executeList(mt.SouthEastAsiaTests)
}

func (T *TEST) Oceania() {
	T.executeList(mt.OceaniaTests)
}

func (T *TEST) AI() {
	T.executeList(mt.AITests)
}


