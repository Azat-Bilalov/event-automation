package storage

type Storage interface {
	Save(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	IsExist(key string) (bool, error)
}
