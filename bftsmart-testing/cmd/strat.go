package cmd

import (
	"github.com/netrixframework/bftsmart-testing/tests"
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/testlib"
)

func GetTest(name string) (*testlib.TestCase, *sm.StateMachine) {
	switch name {
	case "DPropForP":
		return tests.DelayProposeForP(), tests.DelayProposeForPProperty()
	case "DPropSame":
		return tests.DelayProposeSameEpoch(), tests.DelayProposeSameEpochProperty()
	case "DropWrite":
		return tests.DropWrite(), tests.DropWriteProperty()
	case "DropWriteForP":
		return tests.DropWriteForP(), tests.DropWriteForPProperty()
	case "ExpectNewEpoch":
		return tests.ExpectNewEpoch(), tests.ExpectNewEpochProperty()
	case "ExpectStop":
		return tests.ExpectStop(), tests.ExpectStopProperty()
	case "ByzLeaderChange":
		return tests.ByzantineLeaderChange(), tests.ByzantineLeaderChangeProperty()
	case "PrevEpochProposal":
		return tests.PrevEpochProposal(), tests.PrevEpochProposalProperty()
	}
	return nil, nil
}
