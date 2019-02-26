package db

type WriteAction struct {
	DB       *MemDB
	Key      string
	Value    string
	previous string
}

func (a *WriteAction) do() {
	a.previous = a.DB.Read(a.Key)

	a.DB.write(a.Key, a.Value)

	return
}

func (a *WriteAction) rollback() {
	a.DB.write(a.Key, a.previous)
}
