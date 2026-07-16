//
//  ManualBrewRecipe.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation

struct ManualBrewRecipe: Codable, Equatable {
    let method: String
    let coffeeGrams: Double
    let waterGrams: Double
    let grindSize: String
    let temperature: String
    let steps: [String]

    enum CodingKeys: String, CodingKey {
        case method
        case coffeeGrams = "coffee_grams"
        case waterGrams = "water_grams"
        case grindSize = "grind_size"
        case temperature
        case steps
    }
}