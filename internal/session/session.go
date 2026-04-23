package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Message struct {
	Role    string    `json:"role"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

type Session struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	storage  string
}

func NewSessionManager(projectPath string) (*SessionManager, error) {
	storageDir := filepath.Join(projectPath, ".ai-coder-context", "sessions")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create sessions dir: %w", err)
	}

	mgr := &SessionManager{
		sessions: make(map[string]*Session),
		storage:  storageDir,
	}

	// Load existing sessions
	mgr.loadSessions()
	return mgr, nil
}

func (m *SessionManager) CreateSession(projectID string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := &Session{
		ID:        generateID(),
		ProjectID: projectID,
		Messages:  []Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	m.sessions[session.ID] = session
	return session
}

func (m *SessionManager) GetSession(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	return s, ok
}

func (m *SessionManager) AddMessage(sessionID, role, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	session.Messages = append(session.Messages, Message{
		Role:    role,
		Content: content,
		Time:    time.Now(),
	})
	session.UpdatedAt = time.Now()

	return m.saveSession(session)
}

func (m *SessionManager) GetMessages(sessionID string) []Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if session, ok := m.sessions[sessionID]; ok {
		return session.Messages
	}
	return []Message{}
}

func (m *SessionManager) ListSessions() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		result = append(result, s)
	}
	return result
}

func (m *SessionManager) DeleteSession(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[id]; !ok {
		return fmt.Errorf("session not found: %s", id)
	}

	delete(m.sessions, id)

	// Also delete the file
	filePath := filepath.Join(m.storage, id+".json")
	os.Remove(filePath)

	return nil
}

func (m *SessionManager) saveSession(s *Session) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(m.storage, s.ID+".json")
	return os.WriteFile(filePath, data, 0644)
}

func (m *SessionManager) loadSessions() {
	entries, err := os.ReadDir(m.storage)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(m.storage, entry.Name()))
		if err != nil {
			continue
		}

		var session Session
		if err := json.Unmarshal(data, &session); err != nil {
			continue
		}

		m.sessions[session.ID] = &session
	}
}

func generateID() string {
	return fmt.Sprintf("%d-%s", time.Now().Unix(), randomString(8))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}