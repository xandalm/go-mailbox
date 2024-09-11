package memory

import mailbox "github.com/xandalm/go-mailbox"

type provider struct {
	boxes []*box
}

func (p *provider) Create(string) (mailbox.Box, error) {
	b := &box{}
	p.boxes = append(p.boxes, b)
	return b, nil
}

func (p *provider) Delete(string) error {
	panic("unimplemented")
}

func (p *provider) Get(string) (mailbox.Box, error) {
	panic("unimplemented")
}

func (p *provider) List() ([]string, error) {
	panic("unimplemented")
}
