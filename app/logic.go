package app

import (
	"errors"
	"fmt"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrorLocomotiveNotFound = errors.New("locomotive not found")
	ErrorInvalidLocomotive  = errors.New("locomotive is invalid")
	ErrorInternalServer     = errors.New("internal server error")
)

type locomotiveService struct {
	locoRepo LocomotiveRepository
}

func NewLocomotiveService(locoRepo LocomotiveRepository) LocomotiveService {
	return &locomotiveService{
		locoRepo,
	}
}

// Get data from client and store data in datastore
func (l *locomotiveService) CreateLoco(loco *Locomotive) (*Locomotive, error) {

	// Check if for errors from validating locomotive
	if err := validate.Validate(loco); err != nil {

		return nil, errs.Wrap(ErrorInvalidLocomotive, "service.Locomotive.CreateLoco")
	}
	// Define loco id by concatenating shortID twice, just so to make it more unique
	loco.LocoID = fmt.Sprintf("%s%s", shortid.MustGenerate(), shortid.MustGenerate())
	loco.CreatedAT = time.Now().UTC().Unix()
	return l.locoRepo.CreateLoco(loco)
}

// Get loco with loco id
func (l *locomotiveService) GetLoco(loco_id string) (*Locomotive, error) {
	return l.locoRepo.GetLoco(loco_id)
}

// Update loco
func (l *locomotiveService) UpdateLoco(loco *Locomotive) (*Locomotive, error) {
	return l.locoRepo.UpdateLoco(loco)
}

// Delete loco
func (l *locomotiveService) DeleteLoco(loco_id string) error {
	return l.locoRepo.DeleteLoco(loco_id)
}

func (l *locomotiveService) GetAllLoco() (*[]Locomotive, error) {
	return l.locoRepo.GetAllLoco()
}
