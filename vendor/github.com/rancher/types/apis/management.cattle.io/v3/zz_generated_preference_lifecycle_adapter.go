package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type PreferenceLifecycle interface {
	Create(obj *Preference) (runtime.Object, error)
	Remove(obj *Preference) (runtime.Object, error)
	Updated(obj *Preference) (runtime.Object, error)
}

type preferenceLifecycleAdapter struct {
	lifecycle PreferenceLifecycle
}

func (w *preferenceLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*Preference))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *preferenceLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*Preference))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *preferenceLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*Preference))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPreferenceLifecycleAdapter(name string, clusterScoped bool, client PreferenceInterface, l PreferenceLifecycle) PreferenceHandlerFunc {
	adapter := &preferenceLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *Preference) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
