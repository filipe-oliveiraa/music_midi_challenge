package api

type NodeInterface interface {
	Play(bs []byte) error
}
