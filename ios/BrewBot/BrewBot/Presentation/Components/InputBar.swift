//
//  InputBar.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import SwiftUI

struct InputBar: View {
    @Binding var text: String
    let isLoading: Bool
    let onSend: () -> Void

    var body: some View {
        HStack(spacing: 12) {
            // Text input
            TextField("Ask BrewBot...", text: $text, axis: .vertical)
                .lineLimit(1...5)
                .padding(.horizontal, 14)
                .padding(.vertical, 10)
                .background(Color(.secondarySystemBackground), in: Capsule())
                .disabled(isLoading)

            // Send button
            Button(action: onSend) {
                Group {
                    if isLoading {
                        ProgressView()
                            .tint(.white)
                    } else {
                        Image(systemName: "arrow.up")
                            .fontWeight(.semibold)
                    }
                }
                .frame(width: 36, height: 36)
                .background(canSend ? Color.accentColor : Color.secondary, in: Circle())
                .foregroundStyle(.white)
            }
            .disabled(!canSend)
            .animation(.easeInOut(duration: 0.2), value: isLoading)
        }
        .padding(.horizontal)
        .padding(.vertical, 8)
    }

    private var canSend: Bool {
        !text.trimmingCharacters(in: .whitespaces).isEmpty && !isLoading
    }
}
