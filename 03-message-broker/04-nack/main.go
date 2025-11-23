package main

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/gommon/log"
)

type AlarmClient interface {
	StartAlarm() error
	StopAlarm() error
}

func ConsumeMessages(sub message.Subscriber, alarmClient AlarmClient) {
	messages, err := sub.Subscribe(context.Background(), "smoke_sensor")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		var err error
		value := string(msg.Payload)

		switch value {
		case "0":
			err = alarmClient.StopAlarm()
		case "1":
			err = alarmClient.StartAlarm()
		default:
			panic(fmt.Errorf("invalid smoke detection payload: %v", value))
		}

		if err != nil {
			log.Errorf("unable to start/stop alarm: %v", err)
			msg.Nack()
		}
		msg.Ack()
	}
}
