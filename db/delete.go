package db

type DeleteAction struct {
	DB       *MemDB
	Key      string
	previous string
}

func (a *DeleteAction) do() {
	a.previous = a.DB.Read(a.Key)

	a.DB.delete(a.Key)

	return
}

func (a *DeleteAction) rollback() {
	a.DB.write(a.Key, a.previous)
}
