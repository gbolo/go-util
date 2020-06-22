package storage

import "time"

type Service struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Keys        []*ApiKeyStored `json:"api_keys"`
}

func NewService(description string) *Service {
	return &Service{
		ID:          generateID(),
		Description: description,
		Keys:        make([]*ApiKeyStored, 0),
	}
}

func NewServiceWithID(id, description string) *Service {
	return &Service{
		ID:          id,
		Description: description,
		Keys:        make([]*ApiKeyStored, 0),
	}
}

// last used date will be updated if return value is true
func (s *Service) ValidateApiKey(raw string) bool {
	for _, key := range s.Keys {
		if !key.Revoked && key.Validate(raw) {
			key.LastUsed = time.Now().Unix()
			return true
		}
	}
	return false
}

func (s *Service) GenerateApiKey(description string) (rawKey ApiKeyRaw, err error) {
	key, err := GenerateApiKey()
	if err != nil {
		return
	}

	s.addApiKey(&ApiKeyStored{
		Prefix:    key.GetPrefix(),
		Name:      description,
		Hash:      key.GetHash(),
		CreatedAt: time.Now().Unix(),
		Revoked:   false,
	})

	return key, nil
}

func (s *Service) addApiKey(key *ApiKeyStored) {
	s.Keys = append(s.Keys, key)
}

func (s *Service) RevokeApiKey(prefix string) {
	for _, key := range s.Keys {
		if key.Prefix == prefix {
			key.Revoked = true
			key.RevokedAt = time.Now().Unix()
		}
	}
}

func (s *Service) GetApiKey(prefix string) *ApiKeyStored {
	for _, key := range s.Keys {
		if key.Prefix == prefix {
			return key
		}
	}
	return nil
}
