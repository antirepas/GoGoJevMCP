package skill

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	jev "usingSkill/Jev"

	"gopkg.in/yaml.v3"
)

type SkillMetadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
type Skill struct {
	Metadata     SkillMetadata
	Instructions string
}
type SkillRegistry struct {
	Skills map[string]Skill
}

func (r *SkillRegistry) GetSkillInstructions(name string) (string, error) {
	skill, ok := r.Skills[name]
	if !ok {
		return "", errors.New("no such skill exists")
	}
	return skill.Instructions, nil
}

func (r *SkillRegistry) SelectSkill(userInput string) (string, error) {
	jevreq := BuildSkillRequest(userInput, r.Skills)
	resp, err := jev.SendRequest(jevreq)
	if err != nil {
		return "", err
	}
	answer, ok := resp.Answers["skill"]
	if !ok {
		return "", errors.New("Jev did not return a skill answer")
	}
	choice := strings.TrimSpace(answer.Choice)
	if choice == "" {
		return "", errors.New("Jev returned an empty skill choice")
	}
	fmt.Println("Selected skill:", answer.Choice)
	fmt.Println("Confidence:", answer.Confidence)
	fmt.Println("Probabilities:", answer.Probabilities)

	_, ok = r.Skills[choice]
	if !ok && choice != "none" {
		return "", errors.New("this skill does not exist")
	}
	return choice, nil
}

func (r *SkillRegistry) BuildSystemPrompt(skillName string) string {
	var list strings.Builder
	if skillName == "none" {
		list.WriteString("You are a helpful AI assistant, the user is asking about your available skills.\n\nHere are your available skills: \n\n")
		for key := range r.Skills {
			skill := r.Skills[key]
			fmt.Fprintf(&list, "Skill: %s\nDescription: %s\n\n", key, skill.Metadata.Description)
			list.WriteString("\n\n--------------------------------------------\n\n")
		}
		list.WriteString("Explain your capabilities honestly. Do not claim to have skills that are not listed.\n\n")
		return list.String()
	}
	list.WriteString("You are a helpful AI assistant, this is your skill\n\n")
	skill := r.Skills[skillName]
	fmt.Fprintf(&list, "Skill: %s\nDescription: %s\n\n", skillName, skill.Metadata.Description)
	fmt.Fprintf(&list, "Instructions:\n%s\n", skill.Instructions)
	list.WriteString("\n\n--------------------------------------------\n\n")
	return list.String()
}

func BuildSkillRequest(userInput string, skills map[string]Skill) jev.JevRequest {
	criteria := make(map[string]string, len(skills))
	for key, val := range skills {
		criteria[key] = val.Metadata.Description
	}
	criteria["none"] = "If no other skills match the user input"

	req := jev.JevRequest{
		Model: "jev-latest",
		State: fmt.Sprintf("User input: %s", userInput),
		Questions: map[string]jev.JevQuestion{
			"skill": {
				Type:         "choice",
				Instructions: "Choose the single most relevant skill for the user's request. Choose none if no skill is appropriate.",
				Criteria:     criteria,
			},
		},
	}
	return req
}

func LoadSkill(path string) (SkillMetadata, string, error) {
	skill, err := os.ReadFile(path)
	if err != nil {
		return SkillMetadata{}, "", err
	}

	instructions := string(skill)
	if !strings.HasPrefix(instructions, "---") {
		return SkillMetadata{}, "", errors.New("missing YAML frontmatter")
	}

	parts := strings.SplitN(instructions, "---", 3)
	if len(parts) != 3 {
		return SkillMetadata{}, "", errors.New("invalid frontmatter")
	}

	var skillMeta SkillMetadata
	err = yaml.Unmarshal([]byte(strings.TrimSpace(parts[1])), &skillMeta)
	if err != nil {
		return SkillMetadata{}, "", err
	}
	if skillMeta.Name == "" {
		return SkillMetadata{}, "", errors.New("skill metadata does not have a name")
	}
	if skillMeta.Description == "" {
		return SkillMetadata{}, "", errors.New("skill metadata does not have a description")
	}

	instructions = strings.TrimSpace(parts[2])

	return skillMeta, instructions, nil
}

func LoadSkills(path string) (SkillRegistry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return SkillRegistry{}, err
	}
	var keys []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		keys = append(keys, entry.Name())
	}
	sort.Strings(keys)
	skills := make(map[string]Skill, len(entries))
	for _, entry := range keys {
		skillMetadata, instructions, err := LoadSkill(filepath.Join(path, entry, "SKILL.md"))
		if err != nil {
			return SkillRegistry{},
				fmt.Errorf("loading skill %s: %w", entry, err)
		}
		skills[entry] = Skill{
			Metadata:     skillMetadata,
			Instructions: instructions,
		}
	}
	return SkillRegistry{
		Skills: skills,
	}, nil
}
