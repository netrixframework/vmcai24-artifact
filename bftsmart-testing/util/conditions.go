package util

import (
	"strconv"

	"github.com/netrixframework/netrix/log"
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/types"
)

func IsMessageType(t string) sm.Condition {
	return func(e *types.Event, c *sm.Context) bool {
		bftMessage, ok := GetParsedMessage(e, c)
		if !ok {
			return false
		}
		messageID, _ := e.MessageID()
		c.Logger.With(log.LogParams{
			"message_type": bftMessage.TypeString(),
			"expected":     t,
			"message_id":   messageID,
		}).Debug("Checking message type")
		return bftMessage.TypeString() == t
	}
}

func IsPropose() sm.Condition {
	return IsMessageType(ProposeMessageType)
}

func IsWrite() sm.Condition {
	return IsMessageType(WriteMessageType)
}

func IsAccept() sm.Condition {
	return IsMessageType(AcceptMessageType)
}

func IsEpoch(epoch int) sm.Condition {
	return func(e *types.Event, c *sm.Context) bool {
		bftMessage, ok := GetParsedMessage(e, c)
		if !ok {
			return false
		}
		messageID, _ := e.MessageID()
		c.Logger.With(log.LogParams{
			"message_epoch": bftMessage.Epoch,
			"expected":      epoch,
			"message_id":    messageID,
		}).Debug("Checking message epoch")
		return bftMessage.Epoch == epoch
	}
}

func IsView(view int) sm.Condition {
	return func(e *types.Event, c *sm.Context) bool {
		bftMessage, ok := GetParsedMessage(e, c)
		if !ok {
			return false
		}
		return bftMessage.Number == view
	}
}

func IsNewEpoch() sm.Condition {
	return func(e *types.Event, c *sm.Context) bool {
		if !e.IsGeneric() {
			return false
		}
		eType := e.Type.(*types.GenericEventType)
		return eType.T == "NewEpoch"
	}
}

func IsNewEpochOf(epoch int) sm.Condition {
	return func(e *types.Event, c *sm.Context) bool {
		if !e.IsGeneric() {
			return false
		}
		eType := e.Type.(*types.GenericEventType)
		if eType.T != "NewEpoch" {
			return false
		}
		rEpoch, err := strconv.Atoi(eType.Params["epoch"])
		if err != nil {
			return false
		}
		return rEpoch == epoch
	}
}

func IsDecided() sm.Condition {
	return func(e *types.Event, c *sm.Context) bool {
		if !e.IsGeneric() {
			return false
		}
		eType := e.Type.(*types.GenericEventType)
		return eType.T == "Decided"
	}
}

func IsDecidedOf(epoch int) sm.Condition {
	return func(e *types.Event, c *sm.Context) bool {
		if !e.IsGeneric() {
			return false
		}
		eType := e.Type.(*types.GenericEventType)
		if eType.T != "Decided" {
			return false
		}
		rEpoch, err := strconv.Atoi(eType.Params["epoch"])
		if err != nil {
			return false
		}
		return rEpoch == epoch
	}
}
