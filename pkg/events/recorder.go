package events

import (
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Event struct {
	Annotations map[string]string
	EventType   string
	Reason      string
	Message     string
}

type Recorder struct {
	eventRecorder record.EventRecorder
}

func (r *Recorder) Event(obj client.Object, event *Event) {
	r.eventRecorder.AnnotatedEventf(
		obj,
		event.Annotations,
		event.EventType,
		event.Reason,
		event.Message,
	)
}

func NewRecorder(eventRecorder record.EventRecorder) *Recorder {
	return &Recorder{
		eventRecorder: eventRecorder,
	}
}
