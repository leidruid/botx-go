package accounts

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
)

var ErrUnknownBotAccount = errors.New("unknown bot account")

// Storage stores bot accounts and auth tokens.
type Storage struct {
	mu       sync.RWMutex
	accounts []models.BotAccountWithSecret
	tokens   map[uuid.UUID]string
}

func NewStorage(accounts []models.BotAccountWithSecret) *Storage {
	return &Storage{accounts: accounts, tokens: make(map[uuid.UUID]string)}
}

func (s *Storage) GetAccount(botID uuid.UUID) (models.BotAccountWithSecret, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.accounts {
		if a.ID == botID {
			return a, nil
		}
	}
	return models.BotAccountWithSecret{}, ErrUnknownBotAccount
}

func (s *Storage) CTSURL(botID uuid.UUID) (string, error) {
	acc, err := s.GetAccount(botID)
	if err != nil {
		return "", err
	}
	return acc.CTSURL, nil
}

func (s *Storage) SecretKey(botID uuid.UUID) (string, error) {
	acc, err := s.GetAccount(botID)
	if err != nil {
		return "", err
	}
	return acc.SecretKey, nil
}

func (s *Storage) SetToken(botID uuid.UUID, token string) {
	s.mu.Lock()
	s.tokens[botID] = token
	s.mu.Unlock()
}

func (s *Storage) Token(botID uuid.UUID) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tokens[botID]
	return t, ok
}

func (s *Storage) List() []models.BotAccountWithSecret {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]models.BotAccountWithSecret, len(s.accounts))
	copy(out, s.accounts)
	return out
}

// BuildSignature builds HMAC-SHA256 signature of bot_id using secret key.
// BotX expects uppercase hex string.
func (s *Storage) BuildSignature(botID uuid.UUID) (string, error) {
	acc, err := s.GetAccount(botID)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, []byte(acc.SecretKey))
	mac.Write([]byte(botID.String()))
	sum := mac.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(sum)), nil
}
