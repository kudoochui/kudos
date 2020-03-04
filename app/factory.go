package app

type CreateServerFunc func(serverId string) Server

var createMap = map[string]CreateServerFunc{}

func RegisterCreateServerFunc(serverName string, createFunc CreateServerFunc) {
	createMap[serverName] = createFunc
}

func GetCreateServerFunc(serverName string) CreateServerFunc {
	return createMap[serverName]
}
