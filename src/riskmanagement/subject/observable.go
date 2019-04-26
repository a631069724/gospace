package subject

import (
	"sync"
)

type Observable interface {
	AddObserver(Observer)
	DeleteObserver(Observer)
	NotifyObserver()
	SetChanged()
}

func NewObservable() Observable {
	return &ObservableImp{
		observers: make([]Observer, 100, 100),
	}
}

type ObservableImp struct {
	changed   bool
	observers []Observer
	sync.Mutex
}

func (this *ObservableImp) AddObserver(o Observer) {
	this.Lock()
	defer this.Unlock()
	this.observers = append(this.observers, o)
}

func (this *ObservableImp) DeleteObserver(o Observer) {
	this.Lock()
	defer this.Unlock()
	for i, v := range this.observers {
		if v == o {
			this.observers = append(this.observers[:i], this.observers[i+1:]...)
		}
	}
}

func (this *ObservableImp) SetChanged() {
	this.changed = true
}

func (this *ObservableImp) NotifyObserver() {
	if this.changed {
		for _, v := range this.observers {
			v.Update(this)
		}
		this.changed = false
	}
}
