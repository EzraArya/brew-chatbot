//
//  WidgetRenderer.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import SwiftUI

struct WidgetRenderer: View {
    let model: WidgetModel
    
    var body: some View {
        switch model {
        case .manualBrewRecipe(let recipe):
            ManualBrewRecipeCardView(recipe: recipe)
                
        case .unsupported:
            Text("Widget not supported in this version.")
                .italic()
                .foregroundColor(.secondary)
                .padding()
        }
    }
}