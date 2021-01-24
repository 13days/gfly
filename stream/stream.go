package stream

type ContextKey string

type Stream interface {
	Clone() Stream
}
