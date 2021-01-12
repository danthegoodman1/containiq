package notify

import (
	"fmt"
	"github.com/mclenhard/containiq/pkg/setup"
	"github.com/slack-go/slack"
	v1 "k8s.io/api/core/v1"
	"github.com/sirupsen/logrus"
)

var level = map[string]string{
	"Warning":"danger",
	"Normal":"good",
}

func SendSlackEvent(config *setup.Config,event *v1.Event ) {

	slackMessage := fmt.Sprintf("```%s```", event.Message)
	slackTitle := fmt.Sprintf("%s event for resource %s in namespace %s", event.Type,event.InvolvedObject.Kind,event.Namespace)
	attachment := slack.Attachment{
		Text:    "ContainIQ Kubernetes Notification",
		Color: level[event.Type],


		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: slackTitle,
				Value: slackMessage,
			},
		},
	}
	api := slack.New(config.Source.Slack.Key)
	_, _, err  := api.PostMessage(config.Source.Slack.Channel,slack.MsgOptionText("", false),slack.MsgOptionAttachments(attachment))
	if err != nil {
		logrus.Error(err)
		return
	}
}

