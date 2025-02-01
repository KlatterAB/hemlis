package hemlis

import (
	"fmt"
	"sync"
	"time"

	sdk "github.com/bitwarden/sdk-go"
)

type Manager struct {
	client        sdk.BitwardenClientInterface
	cache         map[string]string
	uuidMap       map[string]string
	mu            sync.RWMutex
	cacheDuration time.Duration
	config        Config
}

type Config struct {
	AccessToken    string
	OrganizationID string
	IdentityURL    string
	APIURL         string
	CacheDuration  time.Duration
}

func New(cfg Config) (*Manager, error) {
	client, err := sdk.NewBitwardenClient(&cfg.APIURL, &cfg.IdentityURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create bitwarden client: %w", err)
	}

	err = client.AccessTokenLogin(cfg.AccessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to login with access token: %w", err)
	}

	mgr := &Manager{
		client:        client,
		cache:         make(map[string]string),
		uuidMap:       make(map[string]string),
		mu:            sync.RWMutex{},
		cacheDuration: cfg.CacheDuration,
		config:        cfg,
	}

	if err := mgr.refreshUUIDMap(); err != nil {
		return nil, fmt.Errorf("failed to refresh UUID map: %w", err)
	}

	return mgr, nil
}

func (m *Manager) GetSecretByName(name string) (string, error) {
	// try  to get uuid from map
	m.mu.RLock()
	uuid, exists := m.uuidMap[name]
	m.mu.RUnlock()

	if !exists {
		// uuid not found, refresh the map
		if err := m.refreshUUIDMap(); err != nil {
			return "", fmt.Errorf("failed to refresh UUID map: %w", err)
		}

		// try again after refresh
		m.mu.RLock()
		uuid, exists = m.uuidMap[name]
		m.mu.RUnlock()

		if !exists {
			return "", fmt.Errorf("secret '%s' not found", name)
		}
	}

	return m.GetSecret(uuid)
}

func (m *Manager) GetSecret(uuid string) (string, error) {
	// check if the secret is in the cache
	m.mu.RLock()
	if value, exists := m.cache[uuid]; exists {
		m.mu.RUnlock()
		return value, nil
	}
	m.mu.RUnlock()

	// fetch the secret from Bitwarden
	secret, err := m.client.Secrets().Get(uuid)
	if err != nil {
		return "", fmt.Errorf("failed to get secret '%s': %w", uuid, err)
	}

	// update the cache
	m.mu.Lock()
	m.cache[uuid] = secret.Value
	m.mu.Unlock()

	return secret.Value, nil
}

func (m *Manager) RefreshCache() error {
	// fetch all secrets from Bitwarden
	secrets, err := m.client.Secrets().List(m.config.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	newCache := make(map[string]string, len(secrets.Data))

	for _, secret := range secrets.Data {
		value, err := m.client.Secrets().Get(secret.ID)
		if err != nil {
			return fmt.Errorf("failed to get secret '%s': %w", secret.ID, err)
		}

		newCache[secret.ID] = value.Value
	}

	// update the cache
	m.mu.Lock()
	m.cache = newCache
	m.mu.Unlock()

	return nil
}

func (m *Manager) ClearCache() {
	m.mu.Lock()
	m.cache = make(map[string]string)
	m.mu.Unlock()
}

func (m *Manager) refreshUUIDMap() error {
	secrets, err := m.client.Secrets().List(m.config.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	newMap := make(map[string]string, len(secrets.Data))
	for _, secret := range secrets.Data {
		newMap[secret.Key] = secret.ID
	}

	m.mu.Lock()
	m.uuidMap = newMap
	m.mu.Unlock()

	return nil
}
