package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type textOptions struct {
	// Level: Allowed values: A1, A2, B1, B2, C1, C2
	Level string `json:"level"`

	// Style: Allowed values: formal, informal, friendly, professional, casual, bureaucratic, or other
	Style string `json:"style"`

	// Intention/Tone: Allowed values: aggressive, polite, unsure, confident, apologetic, demanding, grateful, disappointed, or other
	IntentionTone string `json:"intention_tone"`

	// Audience: Allowed values: state office, academic institution, employer, friend, family member, customer service, landlord, teacher, colleague, or other
	Audience string `json:"audience"`

	// Text Type: Allowed values: email, letter, application, complaint, invitation, request, apology, inquiry, announcement, instruction, or other
	TextType string `json:"text_type"`

	// Purpose: Allowed values: to inform, to persuade, to request, to complain, to apologize, to invite, to decline, to confirm, to thank, or other
	Purpose string `json:"purpose"`

	// Context/Situation: Allowed values: job application, apartment search, doctor's appointment, return-exchange, travel booking, university enrollment, insurance claim, neighbor dispute, event planning, or other
	ContextSituation string `json:"context_situation"`

	// Urgency: Allowed values: routine, time-sensitive, urgent, flexible deadline, or other
	Urgency string `json:"urgency"`

	// Relationship to Recipient: Allowed values: stranger, acquaintance, close friend, superior, subordinate, peer, service provider, or other
	RelationshipToRecipient string `json:"relationship_to_recipient"`

	// Complexity: Allowed values: simple single request, multi-part inquiry, complex negotiation, detailed explanation required, or other
	Complexity string `json:"complexity"`

	// Emotional Context: Allowed values: neutral, frustrated, excited, worried, satisfied, angry, hopeful, or other
	EmotionalContext string `json:"emotional_context"`

	// Register: Allowed values: colloquial, standard, elevated, technical-specialized, or other
	Register string `json:"register"`
}

type TextGenerator struct {
	llm *openai.LLM
}

func NewTextGenerator() (*TextGenerator, error) {
	llm, err := openai.New()
	if err != nil {
		return nil, err
	}
	return &TextGenerator{llm: llm}, nil
}

func (tg *TextGenerator) Generate(ctx context.Context, level string) (string, error) {
	opts, err := tg.randomTextOptions(ctx, level)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`Generate a German writing practice text that follows these settings:
- Level: %s
- Style: %s
- Intention/Tone: %s
- Audience: %s
- Text Type: %s
- Purpose: %s
- Context/Situation: %s
- Urgency: %s
- Relationship to Recipient: %s
- Complexity: %s
- Emotional Context: %s
- Register: %s

Only output the text itself without explanations or translations.`, opts.Level, opts.Style, opts.IntentionTone, opts.Audience, opts.TextType, opts.Purpose, opts.ContextSituation, opts.Urgency, opts.RelationshipToRecipient, opts.Complexity, opts.EmotionalContext, opts.Register)

	completion, err := llms.GenerateFromSinglePrompt(ctx, tg.llm, prompt)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(completion), nil
}

// randomTextOptions asks the LLM to pick coherent random values for the text options.
func (tg *TextGenerator) randomTextOptions(ctx context.Context, level string) (*textOptions, error) {
	prompt := fmt.Sprintf(`Create random but coherent settings for a German writing practice task.
Use the CEFR level %s.
Return a JSON object with exactly these keys: level, style, intention_tone, audience, text_type, purpose, context_situation, urgency, relationship_to_recipient, complexity, emotional_context, register.
Choose realistic combinations while still varying the scenario. Only output the JSON.`, level)

	completion, err := llms.GenerateFromSinglePrompt(ctx, tg.llm, prompt)
	if err != nil {
		return nil, err
	}

	var opts textOptions
	if err := json.Unmarshal([]byte(completion), &opts); err != nil {
		return nil, fmt.Errorf("parse text options: %w", err)
	}
	if opts.Level == "" {
		opts.Level = level
	}

	return &opts, nil
}
