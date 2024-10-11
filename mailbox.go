package mailbox

import "time"

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

type Manager interface {
	// Create or restore a box.
	RequestBox(string) (Box, Error)
	// Remove box and all its contents.
	EraseBox(string) Error
	// Check if the box exists
	ContainsBox(string) (bool, Error)
}

type Provider interface {
	// Create a new box.
	Create(string) (Box, Error)
	// Get existing box.
	Get(string) (Box, Error)
	// Check for box existence.
	Contains(string) (bool, Error)
	// Delete existing box and all its contents.
	Delete(string) Error
}

type Bytes []byte

type Data struct {
	CreationTime int64
	Content      Bytes
}

type Box interface {
	// Posts content and return the creation timestamp.
	Post(string, Bytes) (*time.Time, Error)
	// Reads the content matching to the identifier.
	Get(string) (Data, Error)
	// Periodically reads the content matching the given identifiers.
	// Successfully read data will be sent to the channel provided by this method.
	LazyGet(...string) (chan []Data, Error)
	// Lists content identifiers that creation time falls between the given times.
	// The list is sorted by creation time in ascending mode.
	ListFromPeriod(time.Time, time.Time) ([]string, Error)
	// // Lists up to n identifiers of the most recently added content.
	// // The list is sorted by creation time in ascending mode.
	// ListLatest(int64) ([]string, Error)
	// Removes the content matching to the identifier.
	Delete(string) Error
	// Removes all its existing contents.
	Clean() Error
}

type manager struct {
	p Provider
}

func NewManager(p Provider) Manager {
	return &manager{p: p}
}

func (m *manager) RequestBox(id string) (Box, Error) {
	box, err := m.p.Get(id)
	if err != nil {
		return nil, err
	}
	if box == nil {
		box, err = m.p.Create(id)
	}
	return box, err
}

func (m *manager) EraseBox(id string) Error {
	has, err := m.p.Contains(id)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	return m.p.Delete(id)
}

func (m *manager) ContainsBox(id string) (bool, Error) {
	return m.p.Contains(id)
}
