package common

type WS interface {
	SimpleSocket()
	WriteToAllClients()
	WriteToAnUser()
}

type UseCase interface {
	Execute() ([]byte, error)
}
