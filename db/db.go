package db

type RDBMS interface {
	Connect()
	Disconnect()
}
