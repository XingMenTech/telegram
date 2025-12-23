package telegram

const (
	StoreTypeRedis = "redis"
	StoreTypeList  = "list"
)

type Store interface {
	RPush(value string) error
	BLPop() (string, error)
	Size() int64
	Close() error
}

func NewMessageStore(storeType string) Store {
	if StoreTypeRedis == storeType {
		return NewRedisStore("", "", 0)
	} else if StoreTypeList == storeType {
		return NewListStore()
	} else {
		return nil
	}
}
