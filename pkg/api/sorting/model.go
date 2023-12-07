package sorting

type Options struct {
	Fields []string
	Orders []string
}

func NewOptions(fields, orders []string) *Options {
	return &Options{
		Fields: fields,
		Orders: orders,
	}
}
