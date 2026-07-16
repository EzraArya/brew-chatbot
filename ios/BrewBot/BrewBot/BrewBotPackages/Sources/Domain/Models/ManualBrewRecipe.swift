//
//  ManualBrewRecipe.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation

public struct ManualBrewRecipe: Codable, Equatable {
    public let method: String
    public let coffeeGrams: Double
    public let waterGrams: Double
    public let grindSize: String
    public let temperature: String
    public let steps: [String]

    enum CodingKeys: String, CodingKey {
        case method
        case coffeeGrams = "coffee_grams"
        case waterGrams = "water_grams"
        case grindSize = "grind_size"
        case temperature
        case steps
    }
}