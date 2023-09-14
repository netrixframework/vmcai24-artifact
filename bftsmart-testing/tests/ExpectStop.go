package tests

import (
	"time"

	"github.com/netrixframework/bftsmart-testing/util"
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
)

func ExpectStop() *testlib.TestCase {
	filters := testlib.NewFilterSet()

	filters.AddFilter(
		testlib.If(
			util.IsAccept().
				And(util.IsEpoch(0)).
				And(sm.IsMessageTo(types.ReplicaID("3"))),
		).Then(testlib.DropMessage()),
	)

	testCase := testlib.NewTestCase("ExpectStop", 2*time.Minute, ExpectStopProperty(), filters)
	return testCase
}

func ExpectStopProperty() *sm.StateMachine {
	property := sm.NewStateMachine()
	property.Builder().On(
		util.IsMessageType(util.StopMessageType).
			And(sm.IsMessageFrom(types.ReplicaID("3"))),
		sm.SuccessStateLabel,
	)

	return property
}
