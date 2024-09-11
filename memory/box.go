package memory

import "github.com/xandalm/go-mailbox"

type box struct{}

func (b *box) Clean() mailbox.Error {
	panic("unimplemented")
}

func (b *box) Delete(any) mailbox.Error {
	panic("unimplemented")
}

func (b *box) Get(any) (any, mailbox.Error) {
	panic("unimplemented")
}

func (b *box) Post(any, any) mailbox.Error {
	panic("unimplemented")
}
