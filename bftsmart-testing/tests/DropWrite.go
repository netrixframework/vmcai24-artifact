package tests

import (
	"time"

	"github.com/netrixframework/bftsmart-testing/util"
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/testlib"
)

func DropWrite() *testlib.TestCase {
	stateMachine := sm.NewStateMachine()

	filters := testlib.NewFilterSet()

	filters.AddFilter(
		testlib.If(util.IsWrite().And(util.IsEpoch(0))).Then(testlib.DropMessage()),
	)

	testCase := testlib.NewTestCase("DropWrite", 2*time.Minute, stateMachine, filters)
	return testCase
}

func DropWriteProperty() *sm.StateMachine {
	property := sm.NewStateMachine()
	property.Builder().On(
		util.IsNewEpoch(), sm.SuccessStateLabel,
	)
	return property
}
