package session

import (
	"context"
	"errors"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
)

var (
	defaultStore Store
)

type Store interface {
	Get(ctx context.Context) Session
	New(ctx context.Context) (Session, error)
	Delete(ctx context.Context)
}

type Session interface {
	ID() string
	GetInt64(key string) int64
	GetString(key string) string
	GetBool(key string) bool
	Set(key string, value interface{})
}

func newSession(id string) Session {
	return &session{
		id:     id,
		values: make(map[string]interface{}),
	}
}

type session struct {
	id     string
	values map[string]interface{}
}

func (s *session) ID() string {
	return s.id
}

func (s *session) GetInt64(key string) int64 {
	return s.values[key].(int64)
}

func (s *session) GetString(key string) string {
	return s.values[key].(string)
}
func (s *session) GetBool(key string) bool {
	return s.values[key] == true
}

func (s *session) Set(key string, value interface{}) {
	s.values[key] = value
}

func NewStore() Store {
	return &storeImpl{
		sessions: make(map[string]Session),
	}
}

type storeImpl struct {
	sessions map[string]Session
	mu       sync.Mutex
}

func (st *storeImpl) Get(ctx context.Context) Session {
	id := metautils.ExtractIncoming(ctx).Get("x-request-id")
	if id == "" {
		return nil
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	ss := st.sessions[id]
	return ss
}

func (st *storeImpl) Delete(ctx context.Context) {
	id := metautils.ExtractIncoming(ctx).Get("x-request-id")
	if id == "" {
		return
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	delete(st.sessions, id)
}

func (st *storeImpl) New(ctx context.Context) (Session, error) {
	id := metautils.ExtractIncoming(ctx).Get("x-request-id")
	if id == "" {
		return nil, errors.New("x-request-id is missed")
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	ss := newSession(id)
	st.sessions[id] = ss
	return ss, nil
}

func UseDefaultStore(store Store) {
	defaultStore = store
}

func Of(ctx context.Context) Session {
	return defaultStore.Get(ctx)
}
