package tests

import (
	"time"

	"github.com/netrixframework/bftsmart-testing/util"
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
)

func DropWriteForP() *testlib.TestCase {
	stateMachine := sm.NewStateMachine()

	filters := testlib.NewFilterSet()

	filters.AddFilter(
		testlib.If(sm.IsMessageTo(types.ReplicaID("3")).And(util.IsWrite())).
			Then(testlib.DropMessage()),
	)

	testCase := testlib.NewTestCase("DropWriteForP", 2*time.Minute, stateMachine, filters)
	return testCase
}

func DropWriteForPProperty() *sm.StateMachine {
	property := sm.NewStateMachine()
	property.Builder().On(
		util.IsDecided(),
		sm.SuccessStateLabel,
	)
	return property
}
