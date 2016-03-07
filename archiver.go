package podule

type zipper struct{}

type Archiver interface {
	Archive(src, dest string) error
}
