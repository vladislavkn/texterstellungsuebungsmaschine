package ai

type text struct {
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
	
}

func NewTextGenerator() *TextGenerator {
	return &TextGenerator{}
}
