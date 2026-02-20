// Package a05_hallucination provides the hallucination guard agent.
//
// This file defines off-topic and on-topic categories for query classification.
package a05_hallucination

import "strings"

// Category represents a topic category for classification.
type Category struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Patterns    []string `json:"patterns"`
	Keywords    []string `json:"keywords"`
	Confidence  float64  `json:"confidence"`
}

// OffTopicCategories manages categories for off-topic detection.
type OffTopicCategories struct {
	offTopic []Category
	onTopic  []Category
}

// NewOffTopicCategories creates a new category manager with built-in categories.
func NewOffTopicCategories() *OffTopicCategories {
	oc := &OffTopicCategories{
		offTopic: []Category{},
		onTopic:  []Category{},
	}

	// Initialize built-in categories
	oc.initOffTopicCategories()
	oc.initOnTopicCategories()

	return oc
}

// initOffTopicCategories sets up categories that indicate off-topic queries.
func (oc *OffTopicCategories) initOffTopicCategories() {
	oc.offTopic = []Category{
		{
			Name:        "weather",
			Description: "Weather and climate queries",
			Patterns: []string{
				`\bweather\b`,
				`\btemperature\b.*\btoday\b`,
				`\bforecast\b`,
				`\brain\b.*\btomorrow\b`,
				`\bsunny\b`,
				`\bcloudy\b`,
			},
			Keywords:   []string{"weather", "temperature", "forecast", "rain", "sunny", "cloudy", "snow", "wind"},
			Confidence: 0.95,
		},
		{
			Name:        "creative_writing",
			Description: "Creative writing requests",
			Patterns: []string{
				`\bwrite.*poem\b`,
				`\bwrite.*story\b`,
				`\bwrite.*song\b`,
				`\bcompose\b`,
				`\bcreate.*haiku\b`,
				`\bcreative\b`,
			},
			Keywords:   []string{"poem", "story", "song", "haiku", "novel", "fiction", "creative", "compose"},
			Confidence: 0.98,
		},
		{
			Name:        "programming",
			Description: "Programming and coding questions",
			Patterns: []string{
				`\bhow.*implement\b.*\bcode\b`,
				`\bwrite.*function\b`,
				`\bdebug.*code\b`,
				`\balgorithm\b`,
				`\bdata structure\b`,
				`\bsource code\b`,
				`\bprogramming\b`,
			},
			Keywords:   []string{"code", "programming", "algorithm", "debug", "function", "class", "variable", "loop", "array"},
			Confidence: 0.85,
		},
		{
			Name:        "general_knowledge",
			Description: "General knowledge and trivia",
			Patterns: []string{
				`\bwho.*president\b`,
				`\bwhat.*capital\b`,
				`\bhistory\b`,
				`\bgeography\b`,
				`\bpopulation.*of\b`,
			},
			Keywords:   []string{"president", "capital", "country", "history", "geography", "population"},
			Confidence: 0.90,
		},
		{
			Name:        "personal_advice",
			Description: "Personal advice and recommendations",
			Patterns: []string{
				`\bshould i\b`,
				`\bwhat should\b`,
				`\badvice.*me\b`,
				`\brecommend.*restaurant\b`,
				`\bbest.*movie\b`,
			},
			Keywords:   []string{"should", "advice", "recommend", "best movie", "restaurant", "book"},
			Confidence: 0.80,
		},
		{
			Name:        "entertainment",
			Description: "Entertainment and media queries",
			Patterns: []string{
				`\bmovie\b`,
				`\bmusic\b`,
				`\bgame\b`,
				`\bsports\b.*\bscore\b`,
				`\bcelebrity\b`,
			},
			Keywords:   []string{"movie", "music", "game", "sports", "celebrity", "tv show", "netflix"},
			Confidence: 0.92,
		},
		{
			Name:        "math_homework",
			Description: "Math problems and homework",
			Patterns: []string{
				`\bsolve.*equation\b`,
				`\bcalculate.*\d+\s*[+\-*/]\s*\d+`,
				`\bhomework\b`,
				`\bmath.*problem\b`,
				`\bcalculus\b`,
			},
			Keywords:   []string{"equation", "homework", "calculus", "algebra", "geometry"},
			Confidence: 0.88,
		},
		{
			Name:        "travel",
			Description: "Travel and booking queries",
			Patterns: []string{
				`\bflight\b`,
				`\bhotel\b`,
				`\bvacation\b`,
				`\btravel.*to\b`,
				`\bbook.*ticket\b`,
			},
			Keywords:   []string{"flight", "hotel", "vacation", "travel", "booking", "reservation"},
			Confidence: 0.93,
		},
		{
			Name:        "cooking",
			Description: "Cooking and recipes",
			Patterns: []string{
				`\brecipe\b`,
				`\bcook\b`,
				`\bbingredients\b`,
				`\bhow.*make.*food\b`,
			},
			Keywords:   []string{"recipe", "cook", "ingredients", "bake", "food"},
			Confidence: 0.90,
		},
		{
			Name:        "legal_advice",
			Description: "Legal advice queries",
			Patterns: []string{
				`\blegal.*advice\b`,
				`\blawyer\b`,
				`\bcourt\b`,
				`\blawsuit\b`,
			},
			Keywords:   []string{"legal", "lawyer", "court", "lawsuit", "attorney"},
			Confidence: 0.95,
		},
	}
}

// initOnTopicCategories sets up categories that indicate valid healthcare/finance queries.
func (oc *OffTopicCategories) initOnTopicCategories() {
	oc.onTopic = []Category{
		{
			Name:        "healthcare_analytics",
			Description: "Healthcare data and analytics queries",
			Patterns: []string{
				`\bpatient.*visit\b`,
				`\bappointment.*count\b`,
				`\bdoctor.*performance\b`,
				`\bclinic.*metric\b`,
				`\bpharmacy.*sales\b`,
			},
			Keywords: []string{
				"patient", "visit", "appointment", "doctor", "clinic", "hospital",
				"treatment", "diagnosis", "prescription", "pharmacy", "medical",
				"healthcare", "health", "admission", "discharge", "bed",
			},
		},
		{
			Name:        "financial_analytics",
			Description: "Financial and revenue analytics queries",
			Patterns: []string{
				`\brevenue.*by\b`,
				`\btotal.*amount\b`,
				`\bprofit.*margin\b`,
				`\bbilling.*summary\b`,
				`\bpayment.*status\b`,
			},
			Keywords: []string{
				"revenue", "profit", "cost", "expense", "billing", "payment",
				"amount", "price", "income", "earnings", "margin", "budget",
				"financial", "money", "cash", "invoice", "transaction",
			},
		},
		{
			Name:        "operational_metrics",
			Description: "Operational and business metrics",
			Patterns: []string{
				`\bshow.*trend\b`,
				`\bcompare.*performance\b`,
				`\bdepartment.*statistics\b`,
				`\bmonthly.*report\b`,
				`\bkpi\b`,
			},
			Keywords: []string{
				"trend", "compare", "report", "statistics", "metric", "kpi",
				"performance", "growth", "rate", "average", "total", "count",
				"summary", "breakdown", "analysis", "dashboard",
			},
		},
		{
			Name:        "inventory_management",
			Description: "Medicine and supply inventory queries",
			Patterns: []string{
				`\binventory.*level\b`,
				`\bstock.*status\b`,
				`\bmedicine.*available\b`,
				`\bsupply.*chain\b`,
			},
			Keywords: []string{
				"inventory", "stock", "medicine", "supply", "drug", "item",
				"quantity", "available", "order", "reorder", "warehouse",
			},
		},
		{
			Name:        "staff_management",
			Description: "Staff and employee related queries",
			Patterns: []string{
				`\bstaff.*schedule\b`,
				`\bemployee.*count\b`,
				`\bshift.*report\b`,
				`\bpayroll.*summary\b`,
			},
			Keywords: []string{
				"staff", "employee", "doctor", "nurse", "shift", "schedule",
				"payroll", "salary", "working", "hours", "attendance",
			},
		},
	}
}

// GetOffTopicCategories returns all off-topic categories.
func (oc *OffTopicCategories) GetOffTopicCategories() []Category {
	return oc.offTopic
}

// GetOnTopicCategories returns all on-topic categories.
func (oc *OffTopicCategories) GetOnTopicCategories() []Category {
	return oc.onTopic
}

// AddOffTopicCategory adds a new off-topic category.
func (oc *OffTopicCategories) AddOffTopicCategory(category Category) {
	oc.offTopic = append(oc.offTopic, category)
}

// AddOnTopicCategory adds a new on-topic category.
func (oc *OffTopicCategories) AddOnTopicCategory(category Category) {
	oc.onTopic = append(oc.onTopic, category)
}

// IsHealthcareDomain checks if a term is healthcare-specific.
func (oc *OffTopicCategories) IsHealthcareDomain(term string) bool {
	healthcareTerms := map[string]bool{
		"patient":    true,
		"doctor":     true,
		"clinic":     true,
		"hospital":   true,
		"pharmacy":   true,
		"medicine":   true,
		"treatment":  true,
		"diagnosis":  true,
		"prescription": true,
		"appointment": true,
		"visit":      true,
		"admission":  true,
		"discharge":  true,
	}
	return healthcareTerms[strings.ToLower(term)]
}

// IsFinanceDomain checks if a term is finance-specific.
func (oc *OffTopicCategories) IsFinanceDomain(term string) bool {
	financeTerms := map[string]bool{
		"revenue":    true,
		"profit":     true,
		"cost":       true,
		"expense":    true,
		"billing":    true,
		"payment":    true,
		"invoice":    true,
		"transaction": true,
		"budget":     true,
		"margin":     true,
		"earnings":   true,
		"income":     true,
	}
	return financeTerms[strings.ToLower(term)]
}
