//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=cfg.yaml ../../petstore-expanded.yaml

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/uptrace/bunrouter"
)

type PetStore struct {
	Pets   map[int64]Pet
	NextId int64
	Lock   sync.Mutex
}

var _ ServerInterface = (*PetStore)(nil)

func NewPetStore() *PetStore {
	return &PetStore{
		Pets:   make(map[int64]Pet),
		NextId: 1000,
	}
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendPetStoreError(w http.ResponseWriter, code int, message string) {
	petErr := Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(petErr)
}

// Returns all pets
// (GET /pets)
func (p *PetStore) FindPets(w http.ResponseWriter, req bunrouter.Request, params FindPetsParams) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	var result []Pet

	for _, pet := range p.Pets {
		if params.Tags != nil {
			// If we have tags,  filter pets by tag
			for _, t := range *params.Tags {
				if pet.Tag != nil && (*pet.Tag == t) {
					result = append(result, pet)
				}
			}
		} else {
			// Add all pets if we're not filtering
			result = append(result, pet)
		}

		if params.Limit != nil {
			l := int(*params.Limit)
			if len(result) >= l {
				// We're at the limit
				break
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	return bunrouter.JSON(w, result)
}

// Creates a new pet
// (POST /pets)
func (p *PetStore) AddPet(w http.ResponseWriter, req bunrouter.Request) error {
	// We expect a NewPet object in the request body.
	var newPet NewPet
	if err := json.NewDecoder(req.Body).Decode(&newPet); err != nil {
		sendPetStoreError(w, http.StatusBadRequest, "Invalid format for NewPet")
		return err
	}

	// We now have a pet, let's add it to our "database".

	// We're always asynchronous, so lock unsafe operations below
	p.Lock.Lock()
	defer p.Lock.Unlock()

	// We handle pets, not NewPets, which have an additional ID field
	var pet Pet
	pet.Name = newPet.Name
	pet.Tag = newPet.Tag
	pet.Id = p.NextId
	p.NextId = p.NextId + 1

	// Insert into map
	p.Pets[pet.Id] = pet

	// Now, we have to return the NewPet
	return bunrouter.JSON(w, pet)
}

// Deletes a pet by ID
// (DELETE /pets/{id})
func (p *PetStore) DeletePet(w http.ResponseWriter, req bunrouter.Request, id int64) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	pet, found := p.Pets[id]
	if !found {
		sendPetStoreError(w, http.StatusNotFound, "Pet not found")
		return nil
	}

	return bunrouter.JSON(w, pet)
}

// Returns a pet by ID
// (GET /pets/{id})
func (p *PetStore) FindPetByID(w http.ResponseWriter, req bunrouter.Request, id int64) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, found := p.Pets[id]
	if !found {
		sendPetStoreError(w, http.StatusNotFound, fmt.Sprintf("Could not find pet with ID %d", id))
		return nil
	}
	delete(p.Pets, id)

	w.WriteHeader(http.StatusNoContent)
	return nil
}
