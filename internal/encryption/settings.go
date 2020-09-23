package encryption

import "errors"

type Settings struct {
	TurnedOn bool
	Salt     string
}

func (s *Settings) Validate() error {
	if !s.TurnedOn && len(s.Salt) != 32 {
		return errors.New("salt needs to be 32-byte long")
	}
	return nil
}
