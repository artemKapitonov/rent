package service

type Service struct {
	// TorrentClient
	// Huij
}

type TorrentClient interface {
	DownloadToFile(path string) error
}
