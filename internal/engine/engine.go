package engine

import "github.com/livepeer/swarm-chaos/internal/model"

type (
	ChaosEngine struct {
		playgrounds []model.Playground
	}
)

func NewChaosEngine() *ChaosEngine {
	ce := &ChaosEngine{}
	return ce
}

func (ce *ChaosEngine) AddPlayground(playground model.Playground) {
	ce.playgrounds = append(ce.playgrounds, playground)
}

func (ce *ChaosEngine) Entities() ([]model.Entity, error) {
	res := make([]model.Entity, 0)
	for _, driver := range ce.playgrounds {
		entities, err := driver.Entities()
		if err != nil {
			return nil, err
		}
		for _, e := range entities {
			res = append(res, e)
		}
	}
	return res, nil
}

func (ce *ChaosEngine) EntitiesByLabel(key, value string) ([]model.Entity, error) {
	res := make([]model.Entity, 0)
	for _, driver := range ce.playgrounds {
		entities, err := driver.Entities()
		if err != nil {
			return nil, err
		}
		for _, e := range entities {
			labels := e.Labels()
			if labels[key] == value {
				res = append(res, e)
			}
		}
	}
	return res, nil
}
