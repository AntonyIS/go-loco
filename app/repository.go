package app

type LocomotiveRepository interface {
	CreateLoco(loco *Locomotive) (*Locomotive, error)
	GetLoco(loco_id string) (*Locomotive, error)
	UpdateLoco(loco *Locomotive) (*Locomotive, error)
	DeleteLoco(loco_id string) error
	GetAllLoco() (*[]Locomotive, error)
}
