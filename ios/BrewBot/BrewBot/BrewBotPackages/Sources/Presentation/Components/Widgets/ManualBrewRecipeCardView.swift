//
//  ManualBrewRecipeCardView.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import SwiftUI
import Domain

struct ManualBrewRecipeCardView: View {
    let recipe: ManualBrewRecipe
    
    // We store which steps the user has tapped!
    @State private var completedSteps: Set<Int> = []
    
    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            
            // --- HEADER ---
            HStack {
                Image(systemName: "cup.and.saucer.fill")
                    .font(.title2)
                    .foregroundColor(.brown)
                Text(recipe.method)
                    .font(.headline)
                    .bold()
                Spacer()
            }
            
            Divider()
            
            // --- STATS GRID ---
            LazyVGrid(columns: [GridItem(.flexible()), GridItem(.flexible())], spacing: 12) {
                StatItem(icon: "scale.3d", title: "Coffee", value: "\(Int(recipe.coffeeGrams))g")
                StatItem(icon: "drop.fill", title: "Water", value: "\(Int(recipe.waterGrams))g")
                StatItem(icon: "circle.grid.cross", title: "Grind", value: recipe.grindSize)
                StatItem(icon: "thermometer", title: "Temp", value: recipe.temperature)
            }
            .padding(.vertical, 4)
            
            Divider()
            
            // --- INTERACTIVE STEPS ---
            VStack(alignment: .leading, spacing: 14) {
                Text("Instructions")
                    .font(.subheadline)
                    .bold()
                    .foregroundColor(.secondary)
                
                ForEach(Array(recipe.steps.enumerated()), id: \.offset) { index, step in
                    Button(action: {
                        toggleStep(index)
                    }) {
                        HStack(alignment: .top, spacing: 12) {
                            // Dynamic Icon (Circle -> Checkmark)
                            Image(systemName: completedSteps.contains(index) ? "checkmark.circle.fill" : "circle")
                                .foregroundColor(completedSteps.contains(index) ? .green : .gray)
                                .font(.title3)
                            
                            // Dynamic Text (Strikethrough when done)
                            Text(step)
                                .font(.subheadline)
                                .foregroundColor(completedSteps.contains(index) ? .secondary : .primary)
                                .strikethrough(completedSteps.contains(index))
                                .multilineTextAlignment(.leading)
                            
                            Spacer()
                        }
                        .contentShape(Rectangle()) // Makes the whole row tappable
                    }
                    .buttonStyle(.plain)
                }
            }
        }
        .padding(20)
        // Premium card styling
        .background(
            RoundedRectangle(cornerRadius: 20)
                .fill(Color(UIColor.secondarySystemBackground))
                .shadow(color: .black.opacity(0.08), radius: 10, y: 4)
        )
        // Optional: you can constrain the max width so it doesn't stretch too far on iPad
        .frame(maxWidth: 400) 
    }
    
    // Smooth animation when a step is tapped
    private func toggleStep(_ index: Int) {
        withAnimation(.spring(response: 0.3, dampingFraction: 0.6)) {
            if completedSteps.contains(index) {
                completedSteps.remove(index)
            } else {
                completedSteps.insert(index)
            }
        }
    }
}

// A tiny helper view for the 2x2 grid
private struct StatItem: View {
    let icon: String
    let title: String
    let value: String
    
    var body: some View {
        HStack(spacing: 8) {
            Image(systemName: icon)
                .foregroundColor(.brown)
                .frame(width: 24)
            VStack(alignment: .leading, spacing: 2) {
                Text(title)
                    .font(.caption2)
                    .foregroundColor(.secondary)
                    .textCase(.uppercase)
                Text(value)
                    .font(.subheadline)
                    .bold()
            }
            Spacer()
        }
    }
}
