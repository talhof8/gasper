package encryption

import "errors"

type Settings struct {
	TurnedOn bool
	Salt     string
}

func (s *Settings) Validate() error {
	if s.TurnedOn {
		saltBytes := []byte(s.Salt)
		k := len(saltBytes)

		switch k {
		default:
			return errors.New("salt needs to be either 16, 24, or 32-byte long")
		case 16, 24, 32:
			return nil
		}
	}

	return nil
}
