package token

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
)

type TokenOptimizer struct {
	mu       sync.RWMutex
	seenHash map[string]bool
	dedup    bool
}

func NewTokenOptimizer(enableDedup bool) *TokenOptimizer {
	return &TokenOptimizer{
		seenHash: make(map[string]bool),
		dedup:    enableDedup,
	}
}

func (t *TokenOptimizer) Deduplicate(content string) (string, bool) {
	if !t.dedup {
		return content, false
	}

	hash := hashContent(content)
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.seenHash[hash] {
		return "", true // duplicate found
	}

	t.seenHash[hash] = true
	return content, false
}

func (t *TokenOptimizer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.seenHash = make(map[string]bool)
}

func hashContent(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])
}

type CodeCompressor struct{}

func NewCodeCompressor() *CodeCompressor {
	return &CodeCompressor{}
}

func (c *CodeCompressor) Compress(content string, level string) string {
	lines := strings.Split(content, "\n")
	var result []string
	prevEmpty := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip empty lines strategically
		if trimmed == "" {
			if !prevEmpty {
				result = append(result, "")
				prevEmpty = true
			}
			continue
		}
		
		prevEmpty = false

		// Keep non-empty lines
		if strings.HasPrefix(trimmed, "//") || 
		   strings.HasPrefix(trimmed, "#") ||
		   strings.HasPrefix(trimmed, "/*") ||
		   strings.HasPrefix(trimmed, "*/") {
			// Keep comments only if they're meaningful (not just decorations)
			if len(trimmed) > 2 && !strings.Contains(trimmed, "====") && !strings.Contains(trimmed, "----") {
				result = append(result, line)
			}
			continue
		}

		result = append(result, line)
	}

	// Remove trailing empty lines
	for len(result) > 0 && result[len(result)-1] == "" {
		result = result[:len(result)-1]
	}

	return strings.Join(result, "\n")
}

func (c *CodeCompressor) EstimateTokens(text string) int {
	// Rough estimation: ~4 characters per token
	return len(text) / 4
}

func (c *CodeCompressor) OptimizeContext(items []ContextItem, maxTokens int) []ContextItem {
	var result []ContextItem
	currentTokens := 0

	for _, item := range items {
		itemTokens := c.EstimateTokens(item.Content)
		if currentTokens+itemTokens > maxTokens {
			break
		}
		currentTokens += itemTokens
		result = append(result, item)
	}

	return result
}

type ContextItem struct {
	FilePath string
	Content  string
	Type     string
	Priority int
}