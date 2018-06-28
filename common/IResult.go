package common

/*
	Wrapper for a struct
*/
type IResult interface {
	GetChild() IResult
	GetChildren() []chan []string
	WasSuccessful() bool
	Succeed()
	Fail()
	Error() string
	GetMessages() []string
	MergeWithResult(r IResult)
	GetLogLevel() int
	GetStatusCode() int
	SetStatusCode(int)
	GetResponseMessage() string
	SetResponseMessage(string)
	Flush()
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Warningf(template string, args ...interface{})
	DebugMessagef(template string, args ...interface{})
}
