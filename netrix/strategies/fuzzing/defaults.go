package fuzzing

import "github.com/netrixframework/netrix/types"

// Used only for sanity check, says that new states are always there
type DefaultGuider struct {
}

var _ Guider = &DefaultGuider{}

func NewDefaultGuider() *DefaultGuider {
	return &DefaultGuider{}
}

func (d *DefaultGuider) HaveNewState(_ *types.List[*SchedulingChoice], _ *types.List[*types.Event]) bool {
	return true
}

func (d *DefaultGuider) Reset() {}
