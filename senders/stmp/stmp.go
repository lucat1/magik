package stmp

type STMPEmail struct {
}

func (s STMPEmail) Send(auth interface{}, to, body string) error {
	return nil
}

func NewEmail() STMPEmail {
	return STMPEmail{}
}
