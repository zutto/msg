package msg

type middleware interface {
	Wrap(*[]byte) error
	Unwrap(*[]byte) error
}
