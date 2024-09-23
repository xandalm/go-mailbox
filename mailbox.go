package mailbox

import (
	"slices"
	"sync"
)

type Error interface {
	sign() string
	Error() string
}

type mailboxError struct {
	msg string
}

func newError(msg string) *mailboxError {
	return &mailboxError{msg}
}

func (e *mailboxError) sign() string {
	return "mailbox"
}

func (e *mailboxError) Error() string {
	return e.sign() + ": " + e.msg
}

type mailboxDetailedError struct {
	e    *mailboxError
	info string
}

// Returns a detailed error with the addicional info provided
func NewDetailedError(err Error, info string) Error {
	if err, ok := err.(*mailboxError); ok {
		return &mailboxDetailedError{err, info}
	}
	panic("invalid type of error")
}

func (e *mailboxDetailedError) sign() string {
	return e.e.Error()
}

func (e *mailboxDetailedError) Error() string {
	return e.sign() + ", " + e.info
}

var (
	ErrUnableToCreateBox     Error = newError("unable to create the box")
	ErrUnableToRestoreBox    Error = newError("unable to restore box")
	ErrUnableToDeleteBox     Error = newError("unable to delete the box")
	ErrUnableToPostContent   Error = newError("unable to post content")
	ErrUnableToReadContent   Error = newError("unable to read content")
	ErrUnableToDeleteContent Error = newError("unable to delete content")
	ErrUnableToCleanBox      Error = newError("unable to clean box")

	ErrUnknownBox Error = newError("there's no such box")
)

type Provider interface {
	// Create or restore a box.
	RequestBox(string) (Box, Error)
	// Remove box and all its contents.
	EraseBox(string) Error
	// Check if the box exists
	ContainsBox(string) bool
}

type Bytes []byte

type Box interface {
	// Box identifier
	Id() string
	// Post content and return its identifier.
	Post(string, Bytes) Error
	// Read the content matching to the identifier.
	Get(string) (Bytes, Error)
	// Remove the content matching to the identifier.
	Delete(string) Error
	// Remove all its existing contents.
	Clean() Error
}

type Storage interface {
	// Create box container
	CreateBox(string) Error
	// Check box container existence
	List() ([]string, Error)
	// Delete box container
	DeleteBox(string) Error
	// Remove all content from box container
	CleanBox(string) Error
	// Create new content into box container
	CreateContent(string, string, []byte) Error
	// Read content from box container
	ReadContent(string, string) ([]byte, Error)
	// Delete a content from box container
	DeleteContent(string, string) Error
}

type box struct {
	st Storage
	id string
}

func (b *box) Id() string {
	return b.id
}

func (b *box) Clean() Error {
	return b.st.CleanBox(b.id)
}

func (b *box) Delete(id string) Error {
	return b.st.DeleteContent(b.id, id)
}

func (b *box) Get(id string) (Bytes, Error) {
	return b.st.ReadContent(b.id, id)
}

func (b *box) Post(id string, c Bytes) Error {
	return b.st.CreateContent(b.id, id, c)
}

type provider struct {
	mu  sync.Mutex
	st  Storage
	idx []string
}

func NewProvider(st Storage) Provider {
	known, err := st.List()
	slices.Sort(known)
	if err != nil {
		panic("mailbox: unable to load")
	}
	p := &provider{st: st, idx: known}
	return p
}

func (p *provider) contains(id string) (bool, int) {
	pos, ok := slices.BinarySearch(p.idx, id)
	return ok, pos
}

func (p *provider) insert(pos int, id string) {
	p.idx = slices.Insert(p.idx, pos, id)
}

func (p *provider) createBox(pos int, id string) (Box, Error) {
	err := p.st.CreateBox(id)
	if err == nil {
		p.insert(pos, id)
		return &box{id: id}, nil
	}
	return nil, err
}

func (p *provider) RequestBox(id string) (Box, Error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	has, pos := p.contains(id)
	if has {
		return &box{
			st: p.st,
			id: id,
		}, nil
	}
	return p.createBox(pos, id)
}

func (p *provider) EraseBox(id string) Error {
	p.mu.Lock()
	defer p.mu.Unlock()
	has, pos := p.contains(id)
	if !has {
		return ErrUnknownBox
	}
	if err := p.st.DeleteBox(id); err != nil {
		return err
	}
	p.idx = slices.Delete(p.idx, pos, pos+1)
	return nil
}

func (p *provider) ContainsBox(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	has, _ := p.contains(id)
	return has
}
