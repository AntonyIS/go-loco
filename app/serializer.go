package app

type LocomotiveSerializer interface {
	Decode(input []byte) (*Locomotive, error)
	Encode(input *Locomotive) ([]byte, error)
}
