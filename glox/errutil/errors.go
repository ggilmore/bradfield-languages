package errutil

type LoxLanguageError interface {
	IsLoxLanguageError()
	error
}
