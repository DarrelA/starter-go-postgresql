package repository

type RDBMS interface {
	Disconnect()
}

type InMemoryDB interface {
	Disconnect()
}
