package rqueue

type ErrAsyncAPI struct {
	Reason string
}

func (err *ErrAsyncAPI) Error() string {
	return err.Reason
}

var ErrUnexpectedResponseLength = &ErrAsyncAPI{Reason: "unexpected type"}
var ErrUnexpectedMessageType = &ErrAsyncAPI{Reason: "unexpected response length"}
