package db

import (
	"fmt"
	"strings"
)

// Note with this setup, we can't start a transaction with a rollback but I wanted to avoid a memory leak involved in persisting
// a list of all previous actions. This is definitely doable with some more thought as to how far back we're willing to rollback
// and especially with remote storage of action history.
// Concurrent transactions would also greatly confuse that storage.
// Could have added behavior to rollback of previous action to every action (if current action has already rolledback) but
// that would require adding side-effects to the rollback funcs which is a no-no
// Instead we're building the action (for rollback) targeting logic in the Do func

type MemDBTransaction struct {
	Actions []action
}

type action interface {
	do()
	rollback()
}

func (db *MemDB) NewTransaction(inputs []string) (*MemDBTransaction, error) {
	transaction := &MemDBTransaction{
		Actions: []action{},
	}

	rollbackQueue := []action{}

	for _, input := range inputs {
		switch command := strings.Fields(input); command[0] {
		case "WRITE":
			write := db.NewWriteAction(command[1], command[2])
			rollbackQueue = append(rollbackQueue, write)
			transaction.Actions = append(transaction.Actions, write)
		case "DELETE":
			delete := db.NewDeleteAction(command[1])
			rollbackQueue = append(rollbackQueue, delete)
			transaction.Actions = append(transaction.Actions, delete)
		case "ROLLBACK":
			// target last write or delete
			rollbackTarget := rollbackQueue[len(rollbackQueue)-1]
			// remove that action from the queue to rollback
			rollbackQueue = rollbackQueue[:len(rollbackQueue)-1]
			rollback := db.NewRollbackAction(rollbackTarget)
			transaction.Actions = append(transaction.Actions, rollback)
		case "PRINT":
			print := db.NewPrintAction()
			transaction.Actions = append(transaction.Actions, print)
		default:
			return nil, fmt.Errorf("unknown instruction `%s` found", command[0])
		}
	}

	return transaction, nil
}

func (db *MemDB) NewWriteAction(key, value string) *WriteAction {
	return &WriteAction{
		DB:    db,
		Key:   key,
		Value: value,
	}
}

func (db *MemDB) NewDeleteAction(key string) *DeleteAction {
	return &DeleteAction{
		DB:  db,
		Key: key,
	}
}

func (db *MemDB) NewRollbackAction(targetAction action) *RollbackAction {
	return &RollbackAction{
		TargetAction: targetAction,
	}
}

func (db *MemDB) NewPrintAction() *PrintAction {
	return &PrintAction{
		DB: db,
	}
}
