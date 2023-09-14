package util

import (
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
)

var randomProcessKey = "_random_process"

func PickRandomProcess() func(*testlib.Context) error {
	return func(ctx *testlib.Context) error {
		r, ok := ctx.ReplicaStore.GetRandom()
		if ok {
			ctx.Vars.Set(randomProcessKey, string(r.ID))
		}
		return nil
	}
}

func IsMessageFromRandom() sm.Condition {
	return func(e *types.Event, c *sm.Context) bool {
		randomProcess, ok := c.Vars.GetString(randomProcessKey)
		if !ok {
			return false
		}
		message, ok := c.GetMessage(e)
		if !ok {
			return false
		}
		return string(message.From) == randomProcess
	}
}

func IsMessageToRandom() sm.Condition {
	return func(e *types.Event, c *sm.Context) bool {
		randomProcess, ok := c.Vars.GetString(randomProcessKey)
		if !ok {
			return false
		}
		message, ok := c.GetMessage(e)
		if !ok {
			return false
		}
		return string(message.To) == randomProcess
	}
}
