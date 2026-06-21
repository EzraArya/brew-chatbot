//
//  StreamingBubble.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import SwiftUI

struct StreamingBubble: View {
    let content: String

    @State private var cursorOpacity = 1.0

    var body: some View {
        HStack {
            // It always acts like a bot message (left-aligned)
            Text(markdownContent)
                .padding(.horizontal, 14)
                .padding(.vertical, 10)
                .background(Color(.secondarySystemBackground))
                .foregroundStyle(.primary)
                .clipShape(
                    UnevenRoundedRectangle(
                        topLeadingRadius: 4,
                        bottomLeadingRadius: 18,
                        bottomTrailingRadius: 18,
                        topTrailingRadius: 18
                    )
                )

            Spacer(minLength: 60)
        }
        .onAppear {
            // Start the blinking cursor animation
            withAnimation(.easeInOut(duration: 0.6).repeatForever()) {
                cursorOpacity = 0.0
            }
        }
    }

    // MARK: - Helpers
    
    private var markdownContent: AttributedString {
        let lines = content.components(separatedBy: "\n")
        var result = AttributedString()
        
        for (index, line) in lines.enumerated() {
            let parsed = (try? AttributedString(markdown: line)) ?? AttributedString(line)
            result += parsed
            if index < lines.count - 1 {
                result += AttributedString("\n")
            }
        }
        
        // Add the blinking block cursor at the very end
        var cursor = AttributedString(" ▋")
        cursor.foregroundColor = .accentColor.opacity(cursorOpacity)
        result += cursor
        
        return result
    }
}
