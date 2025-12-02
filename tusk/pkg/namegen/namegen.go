package namegen

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	adjectives = []string{
		"akatsuki", "neon", "cosmic", "eternal", "shadow", "crystal", "stellar",
		"azure", "crimson", "emerald", "golden", "silver", "violet", "sakura",
		"mystic", "ancient", "divine", "celestial", "lunar", "solar", "nebula",
		"phantom", "void", "prism", "aurora", "zenith", "nova", "quantum", "infinity",
	}
	nouns = []string{
		"ninja", "samurai", "shinobi", "kitsune", "dragon", "phoenix", "tengu",
		"oni", "yokai", "kami", "tsuki", "hana", "kaze", "mizu", "hi", "kage",
		"ken", "katana", "shuriken", "kunai", "sakura", "moon", "star", "cloud",
		"storm", "wave", "flame", "ice", "thunder", "light",
	}
)

// NameGenerator generates unique anime-themed names
type NameGenerator struct {
	usedNames map[string]bool
	rng       *rand.Rand
}

// NewNameGenerator creates a new name generator instance
func NewNameGenerator() *NameGenerator {
	return &NameGenerator{
		usedNames: make(map[string]bool),
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate creates a unique name in the format: adjective-noun-suffix
func (ng *NameGenerator) Generate() string {
	maxAttempts := 100
	for range maxAttempts {
		adj := adjectives[ng.rng.Intn(len(adjectives))]
		noun := nouns[ng.rng.Intn(len(nouns))]
		suffix := ng.generateSuffix()

		name := fmt.Sprintf("%s-%s-%s", adj, noun, suffix)

		if !ng.usedNames[name] {
			ng.usedNames[name] = true
			return name
		}
	}

	// Fallback: add timestamp if all attempts fail
	adj := adjectives[ng.rng.Intn(len(adjectives))]
	noun := nouns[ng.rng.Intn(len(nouns))]
	suffix := fmt.Sprintf("%s%d", ng.generateSuffix(), time.Now().Unix()%10000)
	name := fmt.Sprintf("%s-%s-%s", adj, noun, suffix)
	ng.usedNames[name] = true
	return name
}

// generateSuffix creates a short unique suffix (4-5 characters)
func (ng *NameGenerator) generateSuffix() string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	length := 4 + ng.rng.Intn(2) // 4 or 5 characters
	suffix := make([]byte, length)
	for i := range suffix {
		suffix[i] = chars[ng.rng.Intn(len(chars))]
	}
	return string(suffix)
}

// GenerateOne creates a single unique name without maintaining state
func GenerateOne() string {
	ng := NewNameGenerator()
	return ng.Generate()
}

func AddSuffix(name string) string {
	ng := NewNameGenerator()
	return fmt.Sprintf("%s-%s", name, ng.generateSuffix())
}
