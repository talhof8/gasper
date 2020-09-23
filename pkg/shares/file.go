package shares

type SharedFile struct {
	ID       string
	Checksum string
	Shares   []*Share
}
