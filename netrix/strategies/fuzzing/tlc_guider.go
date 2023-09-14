package fuzzing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/netrixframework/netrix/types"
)

type TLCGuider struct {
	stateMap map[int64]bool
	tlcAddr  string
	mapper   TLCMapper
}

type TLCMapper func(*types.List[*types.Event]) []TlcEvent

var _ Guider = &TLCGuider{}

func NewTLCGuider(tlcAddr string, mapper TLCMapper) *TLCGuider {
	return &TLCGuider{
		tlcAddr:  tlcAddr,
		stateMap: make(map[int64]bool),
		mapper:   mapper,
	}
}

func (t *TLCGuider) HaveNewState(trace *types.List[*SchedulingChoice], eventTrace *types.List[*types.Event]) bool {
	bs, err := json.Marshal(t.mapEventTrace(eventTrace))
	if err != nil {
		return false
	}
	resp, err := http.Post("http://"+t.tlcAddr+"/execute", "application/json", bytes.NewBuffer(bs))
	if err != nil {
		panic(fmt.Sprintf("failed to communicate with TLC: %s", err))
		return false
	}
	defer resp.Body.Close()
	respS, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	tlcResp := &tlcResponse{}
	err = json.Unmarshal(respS, tlcResp)
	if err != nil {
		return false
	}

	haveNew := false
	for _, k := range tlcResp.Keys {
		if _, ok := t.stateMap[k]; !ok {
			haveNew = true
			t.stateMap[k] = true
		}
	}
	return haveNew
}

func (t *TLCGuider) Reset() {
	t.stateMap = make(map[int64]bool)
}

type TlcEvent struct {
	Name   string
	Params map[string]string
	Reset  bool
}

func (t *TLCGuider) mapEventTrace(events *types.List[*types.Event]) []TlcEvent {
	result := t.mapper(events)
	result = append(result, TlcEvent{Reset: true})
	return result
}

type tlcResponse struct {
	States []string
	Keys   []int64
}

func DefaultEventMapper() TLCMapper {
	return func(l *types.List[*types.Event]) []TlcEvent {
		result := make([]TlcEvent, 0)
		for _, e := range l.Iter() {
			next := TlcEvent{
				Name:   e.TypeS,
				Params: make(map[string]string),
			}
			for k, v := range e.Params {
				next.Params[k] = v
			}
			result = append(result, next)
		}
		return result
	}
}
