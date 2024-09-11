package memory

type box struct{}

func (b *box) Clean() error {
	panic("unimplemented")
}

func (b *box) Delete(any) error {
	panic("unimplemented")
}

func (b *box) Get(any) (any, error) {
	panic("unimplemented")
}

func (b *box) Post(any, any) error {
	panic("unimplemented")
}
