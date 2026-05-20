package deployments

import "slices"

type Store struct {
	deployments []Deployment
}

func NewStore(seed []Deployment) *Store {
	data := slices.Clone(seed)
	return &Store{deployments: data}
}

func (s *Store) List(service, status string) []Deployment {
	filtered := make([]Deployment, 0, len(s.deployments))
	for _, d := range s.deployments {
		if service != "" && d.Service != service {
			continue
		}
		if status != "" && d.Status != status {
			continue
		}
		filtered = append(filtered, d)
	}
	return filtered
}

func (s *Store) GetByID(id string) (Deployment, bool) {
	for _, d := range s.deployments {
		if d.ID == id {
			return d, true
		}
	}
	return Deployment{}, false
}
