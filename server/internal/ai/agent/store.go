package agent

import (
	"context"
	"encoding/base64"
	"mayfly-go/pkg/cache"

	"github.com/cloudwego/eino/compose"
)

type CheckPointStore interface {
	compose.CheckPointStore

	Delete(ctx context.Context, key string) error
}

var checkPointStore CheckPointStore

func GetDefaultCheckPointStore() CheckPointStore {
	if checkPointStore != nil {
		return checkPointStore
	}
	checkPointStore = NewCheckPointStore()
	return checkPointStore
}

func NewCheckPointStore() CheckPointStore {
	return &cacheStore{}
}

type cacheStore struct {
}

var _ CheckPointStore = (*cacheStore)(nil)

func (i *cacheStore) Set(ctx context.Context, key string, value []byte) error {
	cache.Set(key, base64.StdEncoding.EncodeToString(value), -1)
	return nil
}

func (i *cacheStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	encoded := cache.GetStr(key)
	if encoded == "" {
		return nil, false, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, false, err
	}
	return decoded, true, nil
}

func (i *cacheStore) Delete(ctx context.Context, key string) error {
	cache.Del(key)
	return nil
}
