package storage

type Storage interface {
	StoreOperation(walletID, amount string) error
}
