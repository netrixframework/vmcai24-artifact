package tests

import (
	"time"

	"github.com/netrixframework/bftsmart-testing/util"
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
)

func DelayProposeForP() *testlib.TestCase {
	filters := testlib.NewFilterSet()

	filters.AddFilter(
		testlib.If(util.IsWrite().And(util.IsEpoch(0))).Then(testlib.DropMessage()),
	)

	filters.AddFilter(
		testlib.If(util.IsPropose().And(sm.IsMessageTo(types.ReplicaID("3"))).And(util.IsEpoch(0))).
			Then(testlib.StoreInSet(sm.Set("reorderedPropose"))),
	)
	filters.AddFilter(
		testlib.If(
			sm.IsMessageReceive().
				And(sm.IsMessageTo(types.ReplicaID("3"))).
				And(util.IsWrite()).
				And(util.IsEpoch(1)),
		).Then(
			testlib.DeliverAllFromSet(sm.Set("reorderedPropose")),
		),
	)

	testCase := testlib.NewTestCase(
		"DelayProposeForP",
		2*time.Minute,
		sm.NewStateMachine(),
		filters,
	)
	return testCase
}

func DelayProposeForPProperty() *sm.StateMachine {
	property := sm.NewStateMachine()
	start := property.Builder()

	start.On(
		sm.IsMessageReceive().And(sm.IsEventOf(types.ReplicaID("3"))).
			And(util.IsWrite()).
			And(util.IsEpoch(1)),
		"Epoch1ProposeReceived",
	).On(
		sm.IsMessageReceive().And(sm.IsEventOf(types.ReplicaID("3"))).
			And(util.IsPropose()).
			And(util.IsEpoch(0)),
		sm.SuccessStateLabel,
	)

	return property
}
