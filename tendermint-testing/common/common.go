package common

import (
	"math/rand"

	"github.com/netrixframework/netrix/log"
	"github.com/netrixframework/netrix/sm"
	"github.com/netrixframework/netrix/testlib"
	"github.com/netrixframework/netrix/types"
	"github.com/netrixframework/tendermint-testing/util"
)

var (
	DefaultOptions = []SetupOption{partition}

	curRoundKey      = "_curRound"
	commitBlockIDKey = "_commitBlockId"
	randomReplicaKey = "_randomReplica"
)

type SetupOption func(*testlib.Context)

func Setup(sysParams *SystemParams, options ...SetupOption) func(*testlib.Context) error {
	return func(c *testlib.Context) error {
		c.Vars.Set("n", sysParams.N)
		c.Vars.Set("faults", sysParams.F)
		if len(options) == 0 {
			options = append(options, DefaultOptions...)
		}
		for _, o := range options {
			o(c)
		}
		return nil
	}
}

func PickRandomReplica() SetupOption {
	return func(c *testlib.Context) {
		rI := rand.Intn(c.ReplicaStore.Cap())
		var replica types.ReplicaID
		for i, r := range c.ReplicaStore.Iter() {
			replica = r.ID
			if i == rI {
				break
			}
		}
		c.Vars.Set(randomReplicaKey, string(replica))
		c.Logger.With(log.LogParams{
			"randomReplica": replica,
		}).Debug("Picked random replica")
	}
}

func GetRandomReplica(_ *types.Event, c *sm.Context) (types.ReplicaID, bool) {
	rS, ok := c.Vars.GetString(randomReplicaKey)
	return types.ReplicaID(rS), ok
}

func partition(c *testlib.Context) {
	f := int((c.ReplicaStore.Cap() - 1) / 3)
	partitioner := util.NewGenericPartitioner(c.ReplicaStore)
	partition, err := partitioner.CreatePartition(
		[]int{1, f, 2 * f},
		[]string{"h", "faulty", "rest"},
	)
	if err != nil {
		c.Logger.With(log.LogParams{
			"error": err,
		}).Debug("Error creating partition")
		return
	}
	c.Vars.Set("partition", partition)
	c.Logger.With(log.LogParams{
		"partition": partition.String(),
	}).Debug("Partitioned replicas")
}

func GetCurRound(ctx *testlib.Context) (int, bool) {
	return ctx.Vars.GetInt(curRoundKey)
}

func GetCommitBlockID(ctx *testlib.Context) (string, bool) {
	return ctx.Vars.GetString(commitBlockIDKey)
}
