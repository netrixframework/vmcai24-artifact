package tests

import (
	"net/http"
	"time"

	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/strategies"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
	"github.com/netrixframework/raft-testing/tests/util"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

func deleteNode(apiAddr, node string) error {
	req, err := http.NewRequest(http.MethodDelete, "http://"+apiAddr+"/"+node, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	return err
}

func pctSetupFunc(recordSetupFunc func(*strategies.Context)) func(*strategies.Context) {
	return func(ctx *strategies.Context) {
		recordSetupFunc(ctx)

		for _, replica := range ctx.ReplicaStore.Iter() {
			addrI, ok := replica.Info["http_api_addr"]
			if !ok {
				continue
			}
			addrS, ok := addrI.(string)
			if !ok {
				continue
			}
			if replica.ID == types.ReplicaID("1") {
				deleteNode(addrS, "4")
			} else if replica.ID == types.ReplicaID("2") {
				deleteNode(addrS, "3")
			}
		}
	}
}

func ConfChangeBugTest() *testlib.TestCase {
	filters := testlib.NewFilterSet()

	testStateMachine := sm.NewStateMachine()
	oneLeader := testStateMachine.Builder().On(util.IsLeader(types.ReplicaID("1")), "OneLeader")

	removed3 := oneLeader.On(util.IsCommitFor(util.Remove(types.ReplicaID("3"))), "Removed3")
	removed3.MarkSuccess()

	filters.AddFilter(
		testlib.If(
			testStateMachine.InState(sm.StartStateLabel).And(
				util.IsMessageType(raftpb.MsgVote).And(
					sm.IsMessageFrom(types.ReplicaID("1")).Not(),
				).Or(
					util.IsMessageType(raftpb.MsgApp).And(sm.IsMessageFrom(types.ReplicaID("1"))),
				),
			),
		).Then(testlib.DropMessage()),
	)
	filters.AddFilter(
		testlib.If(
			testStateMachine.InState("OneLeader").And(
				sm.IsMessageFrom(types.ReplicaID("1")).Or(sm.IsMessageTo(types.ReplicaID("1"))),
			),
		).Then(
			testlib.DropMessage(),
		),
	)
	filters.AddFilter(
		testlib.If(util.IsMessageType(raftpb.MsgApp).And(sm.IsMessageTo(types.ReplicaID("1")))).Then(testlib.DropMessage()),
	)
	filters.AddFilter(
		testlib.If(
			testStateMachine.InState("OneLeader").And(
				util.IsMessageType(raftpb.MsgApp).And(
					sm.IsMessageTo(types.ReplicaID("3")),
				),
			),
		).Then(testlib.DropMessage()),
	)
	filters.AddFilter(
		testlib.If(
			testStateMachine.InState("OneLeader").And(
				util.IsMessageType(raftpb.MsgVote).And(
					sm.IsMessageFrom(types.ReplicaID("1")).Or(sm.IsMessageFrom(types.ReplicaID("3"))),
				),
			),
		).Then(testlib.DropMessage()),
	)
	filters.AddFilter(
		testlib.If(util.IsLeader(types.ReplicaID("1"))).Then(
			testlib.OnceAction("DeleteNode4", util.DeleteNode(types.ReplicaID("4"), types.ReplicaID("1"))),
		),
	)
	filters.AddFilter(
		testlib.If(
			util.IsLeader(types.ReplicaID("2")).Or(util.IsLeader(types.ReplicaID("4"))),
		).Then(
			testlib.OnceAction("DeleteNode3", util.DeleteNode(types.ReplicaID("3"), types.ReplicaID("2"))),
		),
	)
	filters.AddFilter(
		testlib.If(
			testStateMachine.InState("Removed3").Not().And(
				util.IsMessageType(raftpb.MsgApp).And(sm.IsMessageTo(types.ReplicaID("3"))),
			),
		).Then(testlib.DropMessage()),
	)

	testCase := testlib.NewTestCase("ConfigChange", 10*time.Minute, testStateMachine, filters)
	return testCase
}

func ConfChangeBugProperty() *sm.StateMachine {
	property := sm.NewStateMachine()
	start := property.Builder()

	start.On(
		sm.ConditionWithAction(util.IsNewCommit(), util.RecordCommit()),
		sm.StartStateLabel,
	)
	start.On(
		util.IsDifferentCommit(),
		sm.SuccessStateLabel,
	)
	return property
}
