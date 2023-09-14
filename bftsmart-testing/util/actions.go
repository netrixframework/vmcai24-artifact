package util

import (
	"math/rand"

	"github.com/netrixframework/netrix/log"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
)

func PrintMessage() testlib.Action {
	return func(e *types.Event, ctx *testlib.Context) (out []*types.Message) {
		bftMessage, ok := GetParsedMessage(e, ctx.Context)
		if !ok {
			return
		}
		ctx.Logger.With(log.LogParams{"message": bftMessage.String()}).Info("Observed Message")
		return
	}
}

func GarbleValue() testlib.Action {
	return func(e *types.Event, ctx *testlib.Context) []*types.Message {
		message, ok := ctx.GetMessage(e)
		if !ok {
			return []*types.Message{}
		}
		bftMessage, ok := message.ParsedMessage.(*BFTSmartMessage)
		if !ok {
			return []*types.Message{}
		}
		newValue := make([]byte, len(bftMessage.Value))
		rand.Read(newValue)
		newBftMessage := &BFTSmartMessage{
			Number:    bftMessage.Number,
			Epoch:     bftMessage.Epoch,
			PaxosType: bftMessage.PaxosType,
			Value:     newValue,
			Proof:     bftMessage.Proof,
		}
		newMessageBytes, err := newBftMessage.Marshal()
		if err != nil {
			return []*types.Message{}
		}
		return []*types.Message{ctx.NewMessage(message, newMessageBytes, newBftMessage)}
	}
}
