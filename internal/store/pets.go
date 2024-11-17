package store

import "errors"


// 宠物
type Pet struct {
    Species         string `json:"species"`  // 物种 dog cat bird
    Breed           string `json:"breed"`                  // 品种 中华田园犬、威尔士柯基
    MinWeight int `json:"min_weight,omitempty"`
    MaxWeight int `json:"max_weight,omitempty"`
    AverageWeight int `json:"average_weight,omitempty"`
    Weight int `json:"weight,omitempty"`  // 当前的实际体重
    Description string `json:"description,omitempty"`
    LifeSpan int `json:"lifespan,omitempty"`  // 平均寿命
    Color string `json:"color,omitempty"`
    Age int `json:"age,omitempty"`  // 当前的实际年龄
    AgeEstimated bool `json:"age_estimated,omitempty"`  // 当前的实际年龄, 是否为估算
}

type PetInterface interface {
	SetSpecies(s string) *Pet
	SetBreed(s string) *Pet
	SetMinWeight(s int) *Pet
	SetMaxWeight(s int) *Pet
	SetWeight(s int) *Pet
	SetDescription(s string) *Pet
	SetLifeSpan(s int) *Pet
	SetColor(s string) *Pet
	SetAge(s int) *Pet
	SetAgeEstimated(s bool) *Pet
	Build() (*Pet, error)
}

func NewPetBuilder() PetInterface {
	return &Pet{}
}

// SetSpecies sets the species for our pet, and returns a *Pet.
func (p *Pet) SetSpecies(s string) *Pet {
	p.Species = s
	return p
}

// SetBreed sets the breed for our pet, and returns a *Pet.
func (p *Pet) SetBreed(s string) *Pet {
	p.Breed = s
	return p
}

// SetMinWeight sets the minimum weight for our pet, and returns a *Pet.
func (p *Pet) SetMinWeight(s int) *Pet {
	p.MinWeight = s
	return p
}

// SetMaxWeight sets the maximum weight for our pet, and returns a *Pet.
func (p *Pet) SetMaxWeight(s int) *Pet {
	p.MaxWeight = s
	return p
}

// SetWeight sets the maximum weight for our pet, and returns a *Pet.
func (p *Pet) SetWeight(s int) *Pet {
	p.Weight = s
	return p
}

// SetDescription sets the description for our pet, and returns a *Pet.
func (p *Pet) SetDescription(s string) *Pet {
	p.Description = s
	return p
}

// SetLifespan sets the lifespan for our pet, and returns a *Pet.
func (p *Pet) SetLifeSpan(s int) *Pet {
	p.LifeSpan = s
	return p
}

// SetColor sets the color for our pet, and returns a *Pet.
func (p *Pet) SetColor(s string) *Pet {
	p.Color = s
	return p
}

// SetAge sets the age for our pet, and returns a *Pet.
func (p *Pet) SetAge(s int) *Pet {
	p.Age = s
	return p
}

// SetAgeEstimated sets the age for our pet, and returns a *Pet.
func (p *Pet) SetAgeEstimated(s bool) *Pet {
	p.AgeEstimated = s
	return p
}

// Build uses the various "Set" functions above to build a pet, using the
// fluent interface. The inclusion of this function makes this an example
// of the Builder pattern.
func (p *Pet) Build() (*Pet, error) {
	if p.MinWeight > p.MaxWeight {
		return nil, errors.New("minimum weight must be less than maximum weight")
	}

	p.AverageWeight = (p.MinWeight + p.MaxWeight) / 2

	return p, nil
}
