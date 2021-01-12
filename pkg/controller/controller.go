package controller

import (
	"github.com/containiq/containiq/pkg/notify"
	"github.com/containiq/containiq/pkg/setup"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"os"
	"os/signal"
	"syscall"
	"time"
)
type e struct {

}
type Event interface {
	Send(event *v1.Event, config *setup.Config,currentTime time.Time)
}

func Controller(config *setup.Config) {


	configuration, err := rest.InClusterConfig()
	if err != nil {
		logrus.Error(err)
		panic(err.Error())
	}
	client, err := kubernetes.NewForConfig(configuration)
	if err != nil {
		logrus.Error(err)
		panic(err.Error())
	}

	logrus.Info("Controller Started")
	e := e{}
	WatchEvents(config,client,e)
}

func WatchEvents(config *setup.Config, client *kubernetes.Clientset,event Event) {
	currentTime := time.Now()
	factory := informers.NewSharedInformerFactory(client, time.Duration(30)*time.Minute)

	informer := factory.Core().V1().Events().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}){
			objectEvent := obj.(*v1.Event)
			event.Send(objectEvent,config,currentTime)
		},
	})
	stop := make(chan struct{})
	go informer.Run(stop)
	signalChannel := make(chan os.Signal,1)
	signal.Notify(signalChannel,syscall.SIGINT,syscall.SIGTERM)
	<- signalChannel
	close(stop)

}
func checkForValue( a string, list []string) bool {
	if a == "all" {
		return true
	}
	for _,b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func (e) Send(event *v1.Event, config *setup.Config,currentTime time.Time) {

	timeCreated := event.ObjectMeta.CreationTimestamp.Time
	newEvent := currentTime.Before(timeCreated)

	namespaceCheck := checkForValue(event.Namespace , config.Monitoring.Namespaces.Watch)
	objectCheck := checkForValue(event.InvolvedObject.Kind,config.Monitoring.Resource.Watch)
	typeCheck := checkForValue(event.Type, config.Monitoring.Level.Watch)

	if( namespaceCheck && objectCheck && newEvent && typeCheck) {
		if config.Source.Slack.Enabled == true {
			notify.SendSlackEvent(config,event)
		}
		if config.Source.Webhook.Enabled == true {
			post := notify.PostData{
				Kind: event.Type,
				Namespace: event.Namespace,
				Object: event.InvolvedObject.Kind,
				Message: event.Message,
				Cluster: event.ClusterName,
			}
			w := notify.Webhook{URL: config.Source.Webhook.URL,  Post: post }
			err := w.WebhookPost()
			if err != nil {
				logrus.Error(err)
			}
		}
	}

}
