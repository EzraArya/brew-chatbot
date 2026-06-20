//
//  MessageBubble.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import SwiftUI

struct MessageBubble: View {
    let message: Message

    private var isUser: Bool { message.role == .user }

    var body: some View {
        HStack {
            if isUser { Spacer(minLength: 60) }

            Text(markdownContent)
                .padding(.horizontal, 14)
                .padding(.vertical, 10)
                .background(bubbleBackground)
                .foregroundStyle(isUser ? .white : .primary)
                .clipShape(bubbleShape)

            if !isUser { Spacer(minLength: 60) }
        }
    }

    // MARK: - Styling
    @ViewBuilder
    private var bubbleBackground: some View {
        if isUser {
            Color.accentColor
        } else {
            Color(.secondarySystemBackground)
        }
    }

    private var bubbleShape: some Shape {
        UnevenRoundedRectangle(
            topLeadingRadius: isUser ? 18 : 4,
            bottomLeadingRadius: 18,
            bottomTrailingRadius: isUser ? 4 : 18,
            topTrailingRadius: 18
        )
    }
}

extension MessageBubble {
    // MARK: - Helpers
    private var markdownContent: AttributedString {
        let lines = message.content.components(separatedBy: "\n")
    
        var result = AttributedString()
        for (index, line) in lines.enumerated() {
            let parsed = (try? AttributedString(markdown: line)) ?? AttributedString(line)
            result += parsed
            if index < lines.count - 1 {
                result += AttributedString("\n")
            }
        }
        return result
    }
}