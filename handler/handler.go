// Package handler defines the server along the command handlers.
//
// There are three commands that this server handles:
// set_counter
// get_counter
// add_then_mul
//
// The set_counter is updating the in-memory counter. The counter is updatable by anyone.
// The get_counter returns the value of the counter.
//
// add_then_mul is doing two operations at a time.
// it receives three numbers: x, y, z.
// then adds the x and y. The result is then multiplied by z.
// the result if multiplication returned to the requester.
//
// The add_then_mul depends on the sample extension.
// since the add, mul commands are defined there.
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
		return request.Fail("failed to get parameter: " + err.Error())
	}

	counter = newValue

	return request.Ok(key_value.Empty())
}

func onGetCounter(request message.Request, _ *log.Logger, _ ...*remote.ClientSocket) message.Reply {
	parameters := key_value.Empty()
	parameters.Set("counter", counter)

	return request.Ok(parameters)
}

// given two numbers add them, then multiply it:
//
// (x + y) * z
//
// stack is:
// web-proxy source -> proxy controller -> sample service -> extension -> extension
//
// parameters, err := request.Run("calc", "add", parameters)
// request.Run(clientSocket, "mul", parameters)
// request.Fail()
// request.Ok(parameters)
func onAddThenMul(request message.Request, _ *log.Logger, extensions ...*remote.ClientSocket) message.Reply {
	x, _ := request.Parameters.GetUint64("x")
	y, _ := request.Parameters.GetUint64("y")
	z, _ := request.Parameters.GetUint64("z")

	calcClient := remote.FindClient(extensions, calcExtension)

	addParameters := key_value.Empty().
		Set("x", x).
		Set("y", y)

	request.Next("add", addParameters)
	addReply, err := calcClient.RequestRemoteService(&request)

	if err != nil {
		return request.Fail(fmt.Sprintf("request command '%s' to '%s' extension: %v", "add", calcExtension, err))
	}

	sum, err := addReply.GetUint64("z")
	if err != nil {
		return request.Fail("no 'z':" + err.Error())
	}

	mulParameters := key_value.Empty().
		Set("x", z).
		Set("y", sum)
	request.Next("mul", mulParameters)

	mulReply, err := calcClient.RequestRemoteService(&request)
	if err != nil {
		return request.Fail(fmt.Sprintf("request command '%s' to '%s' extension: %v", "mul", calcExtension, err))
	}

	return request.Ok(mulReply)
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
