package repository

import (
	"encoding/json"
	"fgo-calc-backend/internal/model"
	"os"
	"path/filepath"
	"sort"
)

type Repository struct {
	servants      []model.Servant
	craftEssences []model.CraftEssence
	traits        map[int]string
	ceEffects     map[int]map[int]map[string]model.CeEffect
	dominateMap   map[int]int
}

func NewRepository(dataDir string) (*Repository, error) {
	repo := &Repository{}
	if err := repo.loadData(dataDir); err != nil {
		return nil, err
	}
	repo.precompute()
	return repo, nil
}

func (r *Repository) loadData(dataDir string) error {
	svtData, err := os.ReadFile(filepath.Join(dataDir, "servants.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(svtData, &r.servants); err != nil {
		return err
	}

	for i := range r.servants {
		for key, detail := range r.servants[i].Diff {
			traitSet := make(map[int]struct{}, len(detail.Traits))
			for _, traitId := range detail.Traits {
				traitSet[traitId] = struct{}{}
			}
			detail.TraitSet = traitSet
			r.servants[i].Diff[key] = detail
		}
	}

	ceData, err := os.ReadFile(filepath.Join(dataDir, "ces.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(ceData, &r.craftEssences); err != nil {
		return err
	}

	traitMapData, err := os.ReadFile(filepath.Join(dataDir, "names", "traits.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(traitMapData, &r.traits); err != nil {
		return err
	}

	return nil
}

func (r *Repository) precompute() {
	r.precomputeCeEffects()
	r.buildDominateMap()
}

func (r *Repository) precomputeCeEffects() {
	r.ceEffects = make(map[int]map[int]map[string]model.CeEffect)
	for _, ce := range r.craftEssences {
		r.ceEffects[ce.Id] = make(map[int]map[string]model.CeEffect)
		for _, svt := range r.servants {
			r.ceEffects[ce.Id][svt.Id] = make(map[string]model.CeEffect)
			for diffKey, detail := range svt.Diff {
				percent := 0.0
				direct := 0
				for _, filter := range ce.Filters {
					match := true
					if len(filter.Traits) > 0 {
						for _, tr := range filter.Traits {
							if _, ok := detail.TraitSet[tr]; !ok {
								match = false
								break
							}
						}
					}
					if match {
						if filter.Effect > 0 {
							percent += filter.Effect
						} else {
							direct += int(-filter.Effect)
						}
						break
					}
				}
				r.ceEffects[ce.Id][svt.Id][diffKey] = model.CeEffect{Percent: percent, Direct: direct}
			}
		}
	}
}

func (r *Repository) buildDominateMap() {
	nonFilterCe := []model.CraftEssence{}
	for _, ce := range r.craftEssences {
		if len(ce.Filters) == 1 && ce.Cost == 12 {
			if len(ce.Filters[0].Traits) == 0 {
				nonFilterCe = append(nonFilterCe, ce)
			}
		}
	}
	sort.Slice(nonFilterCe, func(i, j int) bool {
		if nonFilterCe[i].Filters[0].Effect != nonFilterCe[j].Filters[0].Effect {
			return nonFilterCe[i].Filters[0].Effect > nonFilterCe[j].Filters[0].Effect
		}
		return nonFilterCe[i].Id < nonFilterCe[j].Id
	})
	r.dominateMap = make(map[int]int)
	for i := 0; i < len(nonFilterCe)-1; i++ {
		r.dominateMap[nonFilterCe[i+1].Id] = nonFilterCe[i].Id
	}
}

func (r *Repository) GetServants() []model.Servant {
	return r.servants
}

func (r *Repository) GetCraftEssences() []model.CraftEssence {
	return r.craftEssences
}

func (r *Repository) GetTraits() map[int]string {
	return r.traits
}

func (r *Repository) GetCeEffects() map[int]map[int]map[string]model.CeEffect {
	return r.ceEffects
}

func (r *Repository) GetDominateMap() map[int]int {
	return r.dominateMap
}

