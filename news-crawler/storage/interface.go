package storage

// Interface for store news
type Interface interface {
	StoreNews(time uint64, news string) error
	GetLastestNews() map[uint64]string
}
