package db

import "fmt"

type PrintAction struct {
	DB *MemDB
}

func (a *PrintAction) do() {
	a.DB.lock.RLock()
	defer a.DB.lock.RUnlock()

	for key, value := range a.DB.storage {
		fmt.Printf("%s %s\n", key, value)
	}

	return
}

func (a *PrintAction) rollback() {
	return
}
