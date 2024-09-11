package memory

import mailbox "github.com/xandalm/go-mailbox"

type provider struct {
	boxes map[string]*box
}

func (p *provider) contains(id string) bool {
	_, ok := p.boxes[id]
	return ok
}

func (p *provider) Create(id string) (mailbox.Box, error) {
	if p.contains(id) {
		return nil, mailbox.ErrBoxIDDuplicity
	}
	b := &box{}
	p.boxes[id] = b
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
