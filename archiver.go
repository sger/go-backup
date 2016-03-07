package podule

type Archiver interface {
	Archive(src, dest string) error
}
