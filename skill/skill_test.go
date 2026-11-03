package skill

import (
	"testing"

	"github.com/joho/godotenv"
)

func Load() {
	_ = godotenv.Load()
}

func TestSelectSkill(t *testing.T) {
	Load()
	registry := &SkillRegistry{
		Skills: map[string]Skill{
			"coding": {
				Metadata: SkillMetadata{
					Name:        "coding",
					Description: "Helps write and debug code",
				},
				Instructions: "You are a coding assistant.",
			},
			"writing": {
				Metadata: SkillMetadata{
					Name:        "writing",
					Description: "Helps write and edit documents",
				},
				Instructions: "You are a writing assistant.",
			},
		},
	}

	choice, err := registry.SelectSkill("Help me debug this Go function")

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Jev selected: %s", choice)
}
