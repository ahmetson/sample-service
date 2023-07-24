package handler

import (
	"fmt"
	"github.com/ahmetson/common-lib/data_type/key_value"
	"github.com/ahmetson/service-lib/communication/command"
	"github.com/ahmetson/service-lib/communication/message"
	"github.com/ahmetson/service-lib/controller"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/service-lib/remote"
)

var counter uint64 = 0
var calcExtension string

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

// given two numbers add them, then multiply it:
//
// (x + y) * z
func onAddThenMul(request message.Request, _ *log.Logger, extensions ...*remote.ClientSocket) message.Reply {
	x, err := request.Parameters.GetUint64("x")
	if err != nil {
		return message.Fail("request.Parameters: %w" + err.Error())
	}
	y, err := request.Parameters.GetUint64("y")
	if err != nil {
		return message.Fail("request.Parameters: %w" + err.Error())
	}
	z, err := request.Parameters.GetUint64("z")
	if err != nil {
		return message.Fail("request.Parameters: %w" + err.Error())
	}

	calcClient := remote.FindClient(extensions, calcExtension)
	if calcClient == nil {
		return message.Fail("no extension: " + calcExtension)
	}

	addRequest := message.Request{
		Command:    "add",
		Parameters: key_value.Empty().Set("x", x).Set("y", y),
	}
	addReply, err := calcClient.RequestRemoteService(&addRequest)
	if err != nil {
		return message.Fail(fmt.Sprintf("request command '%s' to '%s' extension: %v", "add", calcExtension, err))
	}

	sum, err := addReply.GetUint64("z")
	if err != nil {
		return message.Fail("no 'z':" + err.Error())
	}

	mulRequest := message.Request{
		Command:    "mul",
		Parameters: key_value.Empty().Set("x", z).Set("y", sum),
	}

	mulReply, err := calcClient.RequestRemoteService(&mulRequest)
	if err != nil {
		return message.Fail(fmt.Sprintf("request command '%s' to '%s' extension: %v", "mul", calcExtension, err))
	}

	reply := message.Reply{
		Status:     message.OK,
		Message:    "",
		Parameters: mulReply,
	}

	return reply
}

// NewReplier returns the controller
// todo add the extension
func NewReplier(logger *log.Logger, calcExtensionUrl string) (*controller.Controller, error) {
	calcExtension = calcExtensionUrl

	replier, err := controller.NewReplier(logger)
	if err != nil {
		return nil, err
	}

	replier.RequireExtension(calcExtensionUrl)

	setCounter := command.NewRoute("set_counter", onSetCounter)
	getCounter := command.NewRoute("get_counter", onGetCounter)
	addThenMul := command.NewRoute("add_then_mul", onAddThenMul, calcExtension)

	replier.AddRoute(setCounter)
	replier.AddRoute(getCounter)
	replier.AddRoute(addThenMul)

	return replier, nil
}
