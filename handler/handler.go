package handler

import (
	"github.com/ahmetson/common-lib/data_type/key_value"
	"github.com/ahmetson/service-lib/communication/command"
	"github.com/ahmetson/service-lib/communication/message"
	"github.com/ahmetson/service-lib/controller"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/service-lib/remote"
)

var counter uint64 = 0

func onSetCounter(request message.Request, _ *log.Logger, _ ...*remote.ClientSocket) message.Reply {
	newValue, err := request.Parameters.GetUint64("counter")
	if err != nil {
		return message.Fail("failed to get parameter: " + err.Error())
	}

	counter = newValue

	return message.Reply{
		Status:     message.OK,
		Message:    "",
		Parameters: key_value.Empty(),
	}
}

func onGetCounter(_ message.Request, _ *log.Logger, _ ...*remote.ClientSocket) message.Reply {
	parameters := key_value.Empty()
	parameters.Set("counter", counter)
	return message.Reply{
		Status:     message.OK,
		Message:    "",
		Parameters: parameters,
	}
}

// NewReplier returns the controller
// todo add the extension
func NewReplier(logger *log.Logger) (*controller.Controller, error) {
	replier, err := controller.NewReplier(logger)
	if err != nil {
		return nil, err
	}

	setCounter := command.NewRoute("set_counter", onSetCounter)
	getCounter := command.NewRoute("get_counter", onGetCounter)

	replier.AddRoute(setCounter)
	replier.AddRoute(getCounter)

	return replier, nil
}
