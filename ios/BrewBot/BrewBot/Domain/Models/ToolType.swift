//
//  ToolType.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation

enum ToolType: String, Codable, Equatable, Sendable, Hashable {
    case generateBrewRecipe = "generate_brew_recipe"
    case unknown

    init(from decoder: any Decoder) throws {
        let container = try decoder.singleValueContainer()
        let rawValue = try container.decode(String.self)
        self = ToolType(rawValue: rawValue) ?? .unknown
    }
}