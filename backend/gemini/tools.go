package gemini

import (
	"google.golang.org/genai"
)

func GetToolConfig() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{
					buildBrewRecipeTool(),
				},
			},
		},
	}
}

func buildBrewRecipeTool() *genai.FunctionDeclaration {
	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"method": {
				Type: genai.TypeString,
				Description: "e.g., V60, Chemex, Aeropress, French Press",
			},
			"coffee_grams": {
				Type: genai.TypeNumber,
				Description: "Amount of coffee in grams",
			},
			"water_grams": {
				Type: genai.TypeNumber,
				Description: "Amount of water in grams",
			},
			"grind_size": {
				Type: genai.TypeString,
				Description: "e.g., Medium-Fine, Coarse, Fine",
			},
			"temperature": {
				Type: genai.TypeString,
				Description: "Water temperature, e.g., 93°C or 200°F",
			},
			"steps": {
				Type: genai.TypeArray,
				Description: "Sequential brewing steps",
				Items: &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{
			"method", "coffee_grams", "water_grams", "grind_size", "temperature", "steps",
		},
	}

	return &genai.FunctionDeclaration{
		Name: "generate_brew_recipe",
		Description: "Generates a structured manual coffee brewing recipe",
		Parameters: schema,
	}
}