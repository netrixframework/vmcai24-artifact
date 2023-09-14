package tests

import (
	"time"

	"github.com/netrixframework/bftsmart-testing/util"
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
)

func ByzantineLeaderChange() *testlib.TestCase {
	stateMachine := sm.NewStateMachine()

	filters := testlib.NewFilterSet()

	filters.AddFilter(
		testlib.If(sm.IsMessageFrom(types.ReplicaID("1")).And(util.IsEpoch(0)).And(util.IsAccept().Or(util.IsWrite()))).
			Then(testlib.StoreInSet(sm.Set("delayedMessages"))),
	)

	filters.AddFilter(
		testlib.If(util.IsNewEpoch()).
			Then(testlib.DeliverAllFromSet(sm.Set("delayedMessages"))),
	)

	filters.AddFilter(
		testlib.If(sm.IsMessageFrom(types.ReplicaID("3")).And(util.IsEpoch(0)).And(util.IsAccept().Or(util.IsWrite()))).
			Then(util.GarbleValue()),
	)

	testCase := testlib.NewTestCase("ByzantineLeaderChange", 2*time.Minute, stateMachine, filters)
	testCase.SetupFunc(util.PickRandomProcess())
	return testCase
}

func ByzantineLeaderChangeProperty() *sm.StateMachine {
	property := sm.NewStateMachine()

	property.Builder().On(
		util.IsNewEpochOf(1),
		sm.SuccessStateLabel,
	)

	return property
}
