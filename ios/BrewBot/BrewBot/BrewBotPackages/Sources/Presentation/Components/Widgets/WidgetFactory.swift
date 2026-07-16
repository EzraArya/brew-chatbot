//
//  WidgetFactory.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation
import Domain
import Core

enum WidgetModel: Equatable {
    case manualBrewRecipe(ManualBrewRecipe)
    case unsupported
}

struct WidgetFactory {
    static func decode(type: ToolType, payload: String) -> WidgetModel {
        guard let data = payload.data(using: .utf8) else {
            return .unsupported
        }

        switch type {
            case .generateBrewRecipe:
                if let recipe = payload.decode(as: ManualBrewRecipe.self) {
                    return .manualBrewRecipe(recipe)
                }
            case .unknown:
                return .unsupported
        }

        return .unsupported
    }
}
