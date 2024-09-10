package mailbox

type Box struct {
	Id string
}

func CreateBox(id string) (Box, error) {
	return Box{id}, nil
}
