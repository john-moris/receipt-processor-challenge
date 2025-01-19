package db

import (
	"sync"

	"github.com/john-moris/receipt-processor-challenge/internal/domain/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const stillInProgressScore = -1

type Memory struct {
	receipts map[string]int
	lock     sync.RWMutex
	logger   *zap.Logger
}

func NewMemory(logger *zap.Logger) repository.Repository {
	return &Memory{
		receipts: make(map[string]int),
		lock:     sync.RWMutex{},
		logger:   logger,
	}
}

func (m *Memory) Get(id string) (int, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	score, ok := m.receipts[id]
	if !ok {
		return 0, repository.ErrItemNotFound
	}

	m.logger.Info("reading from in-memory repository", zap.Int("score", score), zap.String("id", id))

	if score == -1 {
		return 0, repository.ErrItemStillInProgress
	}

	return score, nil
}

func (m *Memory) Start() string {
	m.lock.Lock()
	defer m.lock.Unlock()

	id := uuid.New().String()

	m.logger.Info("writing into in-memory repository", zap.Int("score", stillInProgressScore), zap.String("id", id))

	m.receipts[id] = stillInProgressScore

	return id
}

func (m *Memory) Finish(id string, score int) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.logger.Info("writing into in-memory repository", zap.Int("score", score), zap.String("id", id))

	m.receipts[id] = score
}
