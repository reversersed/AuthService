package storage

import (
	"time"
)

func (s *storage) CreateNewRefreshPassword(uuid string, refreshpassword []byte, creation time.Time) error {
	return nil
}
func (s *storage) GetFreeRefreshToken(id string, createdTime time.Time) (string, []byte, error) {
	return "", nil, nil
}
func (s *storage) RevokeRefreshToken(rowId string) error {
	return nil
}
