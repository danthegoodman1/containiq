package e2e

import (
	"context"
	"github.com/containiq/containiq/pkg/controller"
	"github.com/containiq/containiq/pkg/setup"
	"github.com/containiq/containiq/test/e2e/config"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"reflect"
	"testing"
	"time"
)
type F struct {

}

func (F) Send(event *v1.Event, config *setup.Config,currentTime time.Time) {
	if event.InvolvedObject.Name == "containiq.com" {
		e := *event
		eventList = append(eventList, e)
	}
}

var eventList []v1.Event

func runE2ETests( t *testing.T) {
	//set config files to current directory
	os.Setenv("NOTIFICATION_FILE_LOCATION", "notification-config.yaml")
	os.Setenv("CONFIG_FILE_LOCATION", "config.yaml")

	framework , err := config.NewFramework()
	if err != nil {
		t.Errorf("Failed %v" , err)
	}
	E := F{}
	go controller.WatchEvents(framework.Config,framework.Client,E)

	testPod := &v1.Pod{
		TypeMeta:   meta_v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:"containiq.com",
			Labels: map[string]string{
				"containiqtest":"true",
			},
		},
		Spec:       v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},

		},
		Status:     v1.PodStatus{},
	}

	_ ,err = framework.Client.CoreV1().Pods("default").Create(context.TODO(), testPod, meta_v1.CreateOptions{})

	if err != nil {
		t.Errorf("error creating pod: %v ",err )
	}
	options := meta_v1.ListOptions{
		FieldSelector:"involvedObject.name=containiq.com",
	}
	ev, err := framework.Client.CoreV1().Events("default").List(context.TODO(),options )
	if err !=nil {
		t.Errorf("Error getting events: %v", err )
	}

	for i,p := range ev.Items {
		reflect.DeepEqual(eventList[i], p)
		if !reflect.DeepEqual(eventList[i],p){
			t.Errorf("Failed wanted %v got %v",eventList[i],p)
		}

	}






}