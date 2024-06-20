package db

type RDBMS interface {
	Connect()
	Disconnect()
}

type InMemoryDB interface {
	Connect()
	Disconnect()
}
