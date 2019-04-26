package Manager

type Customer struct {
	markets []string
}

func (c *Customer) SetMarkets(name []string) {
	c.markets = name
}

func (c *Customer) GetMarkets() []string {
	return c.markets
}

func NewCustomer() *Customer {
	return &Customer{
		markets: []string{"cu1810", "cu1811", "cu1812"},
	}
}
