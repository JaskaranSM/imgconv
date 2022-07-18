package storage

type StorageRepo interface {
	GetImageBytesById(string) ([]byte, error)
	SaveImageBytesId(string, []byte) error
}
