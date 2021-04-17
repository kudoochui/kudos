package protocol

type ISession interface {
	GetNodeId() string
	GetSessionId() int64
	GetUserId() int64
	Get(key string) string
	Set(key, value string)
	GetCache(key string) string
	SetCache(key, value string)
	RemoveCache(key string)
}

type DummySession struct {

}

func NewDummySession() *DummySession {
	return &DummySession{}
}

func (d *DummySession) GetNodeId() string {
	return ""
}

func (d *DummySession) GetSessionId() int64 {
	return 0
}

func (d *DummySession) GetUserId() int64 {
	return 0
}

func (d *DummySession) Get(key string) string {
	return ""
}

func (d *DummySession) Set(key, value string) {

}

func (d *DummySession) GetCache(key string) string {
	return ""
}

func (d *DummySession) SetCache(key, value string) {

}

func (d *DummySession) RemoveCache(key string) {

}