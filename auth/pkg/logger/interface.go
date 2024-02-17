package logger

type UC interface {
	Log(level int64, message string)
}
