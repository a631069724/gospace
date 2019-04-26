package rule

type Ruler interface {
	GetType() interface{}
	Judge(Params) (interface{}, error)
}

type RuleHandler func(Params) (interface{}, error)

type MyRuler struct {
	_type   string
	handler RuleHandler
}

func (m *MyRuler) GetType() interface{} {
	return m._type
}

func (m *MyRuler) Judge(p Params) (interface{}, error) {
	return m.handler(p)
}

func NewMyRule(_type string, proc RuleHandler) Ruler {
	return &MyRuler{
		_type:   _type,
		handler: proc,
	}
}
