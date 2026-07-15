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
					buildBeerRecipeTool(),
					buildTeaRecipeTool(),
					buildKombuchaRecipeTool(),
					buildTroubleshootingTool(),
					buildBrewTimerTool(),
				},
			},
		},
	}
}

// buildBrewRecipeTool generates structured manual coffee brewing recipes (V60, Chemex, etc.)
func buildBrewRecipeTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "generate_brew_recipe",
		Description: "Generates a structured manual coffee brewing recipe. Use for pour-over, immersion, and espresso methods.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"method": {
					Type:        genai.TypeString,
					Description: "Brewing method",
					Enum:        []string{"V60", "Chemex", "AeroPress", "French Press", "Kalita Wave", "Moka Pot", "Cold Brew", "Espresso"},
				},
				"coffee_grams": {
					Type:        genai.TypeNumber,
					Description: "Amount of coffee in grams",
				},
				"water_grams": {
					Type:        genai.TypeNumber,
					Description: "Amount of water in grams",
				},
				"grind_size": {
					Type:        genai.TypeString,
					Description: "Grind size description",
					Enum:        []string{"Extra Fine", "Fine", "Medium-Fine", "Medium", "Medium-Coarse", "Coarse", "Extra Coarse"},
				},
				"temperature": {
					Type:        genai.TypeString,
					Description: "Water temperature, e.g., 93°C or 200°F",
				},
				"steps": {
					Type:        genai.TypeArray,
					Description: "Sequential brewing steps in order",
					Items:       &genai.Schema{Type: genai.TypeString},
				},
			},
			Required:         []string{"method", "coffee_grams", "water_grams", "grind_size", "temperature", "steps"},
			PropertyOrdering: []string{"method", "coffee_grams", "water_grams", "grind_size", "temperature", "steps"},
		},
	}
}

// buildBeerRecipeTool generates structured homebrewing beer recipes.
func buildBeerRecipeTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "generate_beer_recipe",
		Description: "Generates a structured homebrewing beer recipe. Use when the user asks about brewing beer, ales, lagers, or any beer style.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"style": {
					Type:        genai.TypeString,
					Description: "Beer style",
					Enum:        []string{"IPA", "Hazy IPA", "Stout", "Porter", "Wheat Beer", "Pale Ale", "Lager", "Pilsner", "Sour", "Saison", "Barleywine", "Brown Ale", "Red Ale"},
				},
				"batch_size_liters": {
					Type:        genai.TypeNumber,
					Description: "Batch size in liters",
				},
				"yeast": {
					Type:        genai.TypeString,
					Description: "Yeast strain name and/or description",
				},
				"target_abv": {
					Type:        genai.TypeNumber,
					Description: "Target alcohol by volume as a percentage, e.g., 5.5",
				},
				"target_ibu": {
					Type:        genai.TypeNumber,
					Description: "Target bitterness in IBU",
				},
				"fermentation_temp": {
					Type:        genai.TypeString,
					Description: "Fermentation temperature, e.g., 20°C",
				},
				"steps": {
					Type:        genai.TypeArray,
					Description: "Brewing steps from mash to packaging",
					Items:       &genai.Schema{Type: genai.TypeString},
				},
			},
			Required:         []string{"style", "batch_size_liters", "yeast", "target_abv", "steps"},
			PropertyOrdering: []string{"style", "batch_size_liters", "yeast", "target_abv", "target_ibu", "fermentation_temp", "steps"},
		},
	}
}

// buildTeaRecipeTool generates structured tea brewing recipes.
func buildTeaRecipeTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "generate_tea_recipe",
		Description: "Generates a structured tea brewing recipe. Use when the user asks how to brew a specific type of tea.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"tea_type": {
					Type:        genai.TypeString,
					Description: "Type of tea",
					Enum:        []string{"Green", "Black", "Oolong", "White", "Pu-erh", "Herbal", "Matcha", "Chai"},
				},
				"tea_grams": {
					Type:        genai.TypeNumber,
					Description: "Amount of tea leaves in grams",
				},
				"water_ml": {
					Type:        genai.TypeNumber,
					Description: "Amount of water in milliliters",
				},
				"temperature": {
					Type:        genai.TypeString,
					Description: "Water temperature, e.g., 80°C or 175°F",
				},
				"steep_time": {
					Type:        genai.TypeString,
					Description: "Steeping duration, e.g., 3 minutes",
				},
				"vessel": {
					Type:        genai.TypeString,
					Description: "Brewing vessel",
					Enum:        []string{"Teapot", "Gaiwan", "French Press", "Infuser", "Kyusu", "Mug"},
				},
				"steps": {
					Type:        genai.TypeArray,
					Description: "Sequential brewing steps",
					Items:       &genai.Schema{Type: genai.TypeString},
				},
			},
			Required:         []string{"tea_type", "tea_grams", "water_ml", "temperature", "steep_time", "vessel", "steps"},
			PropertyOrdering: []string{"tea_type", "tea_grams", "water_ml", "temperature", "steep_time", "vessel", "steps"},
		},
	}
}

// buildKombuchaRecipeTool generates structured kombucha fermentation recipes.
func buildKombuchaRecipeTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "generate_kombucha_recipe",
		Description: "Generates a structured kombucha fermentation recipe. Use when the user asks about brewing or fermenting kombucha.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"tea_base": {
					Type:        genai.TypeString,
					Description: "Tea base used for the kombucha",
					Enum:        []string{"Black Tea", "Green Tea", "Oolong Tea", "White Tea", "Herbal Tea", "Blend"},
				},
				"sugar_grams": {
					Type:        genai.TypeNumber,
					Description: "Amount of sugar in grams",
				},
				"water_ml": {
					Type:        genai.TypeNumber,
					Description: "Amount of water in milliliters",
				},
				"first_fermentation_days": {
					Type:        genai.TypeNumber,
					Description: "Duration of first fermentation in days",
				},
				"second_fermentation_days": {
					Type:        genai.TypeNumber,
					Description: "Duration of second fermentation in days (for carbonation and flavoring)",
				},
				"steps": {
					Type:        genai.TypeArray,
					Description: "Sequential fermentation steps from brewing to bottling",
					Items:       &genai.Schema{Type: genai.TypeString},
				},
			},
			Required:         []string{"tea_base", "sugar_grams", "water_ml", "first_fermentation_days", "second_fermentation_days", "steps"},
			PropertyOrdering: []string{"tea_base", "sugar_grams", "water_ml", "first_fermentation_days", "second_fermentation_days", "steps"},
		},
	}
}

// buildTroubleshootingTool generates structured brewing troubleshooting guides.
func buildTroubleshootingTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "generate_troubleshooting",
		Description: "Generates a structured troubleshooting guide for brewing problems. Use when the user describes a problem with their brew (e.g., too bitter, too sour, no carbonation).",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"problem": {
					Type:        genai.TypeString,
					Description: "A concise description of the problem, e.g., 'Coffee tastes too bitter'",
				},
				"category": {
					Type:        genai.TypeString,
					Description: "Category of the brewing problem",
					Enum:        []string{"Coffee", "Beer", "Tea", "Kombucha", "General"},
				},
				"possible_causes": {
					Type:        genai.TypeArray,
					Description: "List of possible causes for the problem",
					Items:       &genai.Schema{Type: genai.TypeString},
				},
				"quick_tips": {
					Type:        genai.TypeArray,
					Description: "Actionable quick tips to fix the problem",
					Items:       &genai.Schema{Type: genai.TypeString},
				},
			},
			Required:         []string{"problem", "category", "possible_causes", "quick_tips"},
			PropertyOrdering: []string{"problem", "category", "possible_causes", "quick_tips"},
		},
	}
}

// buildBrewTimerTool generates step-by-step brew timers for real-time guidance.
func buildBrewTimerTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "generate_brew_timer",
		Description: "Generates a step-by-step brew timer for real-time brewing guidance. Use when the user wants to brew RIGHT NOW and needs a countdown timer — distinct from a recipe which is reference material.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"title": {
					Type:        genai.TypeString,
					Description: "Title of the timer, e.g., 'V60 Pour-Over'",
				},
				"method": {
					Type:        genai.TypeString,
					Description: "Brewing method this timer is for",
				},
				"total_duration_seconds": {
					Type:        genai.TypeNumber,
					Description: "Total duration of the brew in seconds",
				},
				"steps": {
					Type:        genai.TypeArray,
					Description: "Ordered list of step instructions to display during the countdown",
					Items:       &genai.Schema{Type: genai.TypeString},
				},
			},
			Required:         []string{"title", "method", "total_duration_seconds", "steps"},
			PropertyOrdering: []string{"title", "method", "total_duration_seconds", "steps"},
		},
	}
}