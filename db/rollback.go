package db

type RollbackAction struct {
	TargetAction action
}

func (a *RollbackAction) do() {
	a.TargetAction.rollback()

	return
}

// no-op
func (a *RollbackAction) rollback() {
	return
}
