package rule

type Params interface {
	SetParam(string, interface{})
	GetParam(string) interface{}
}

type MyParams map[string]interface{}

func NewParams() Params {
	mp := make(MyParams)
	return mp
}

func (p MyParams) SetParam(key string, value interface{}) {
	p[key] = value
}

func (p MyParams) GetParam(key string) interface{} {
	return p[key]
}
