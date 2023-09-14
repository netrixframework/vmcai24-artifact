package tests

import (
	"time"

	"github.com/netrixframework/bftsmart-testing/util"
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
)

func DelayProposeSameEpoch() *testlib.TestCase {
	stateMachine := sm.NewStateMachine()

	filters := testlib.NewFilterSet()

	filters.AddFilter(
		testlib.If(util.IsPropose().And(sm.IsMessageTo(types.ReplicaID("3")))).Then(testlib.StoreInSet(sm.Set("delayedPropose")),
			testlib.DropMessage()),
	)

	filters.AddFilter(
		testlib.If(sm.IsMessageReceive().And(sm.IsMessageTo(types.ReplicaID("3")).And(util.IsAccept()))).
			Then(testlib.DeliverAllFromSet(sm.Set("delayedPropose"))),
	)

	testCase := testlib.NewTestCase("DelayProposeSameEpoch", 2*time.Minute, stateMachine, filters)
	return testCase
}

func DelayProposeSameEpochProperty() *sm.StateMachine {
	property := sm.NewStateMachine()
	property.Builder().On(
		sm.IsEventOf(types.ReplicaID("3")).And(util.IsDecided()),
		sm.SuccessStateLabel,
	)
	return property
}
