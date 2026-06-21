//
//  StreamingBubble.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import SwiftUI

struct StreamingBubble: View {
    let content: String
    @State private var cursorVisible = true

    var body: some View {
        HStack {
            HStack(spacing: 2) {
                Text(markdownContent)
                    .animation(nil, value: content)

                Text("▋")
                    .foregroundStyle(.primary)
                    .opacity(cursorVisible ? 1 : 0)
                    .animation(.easeInOut(duration: 0.6).repeatForever(autoreverses: true), value: cursorVisible)
            }
            .padding(.horizontal, 14)
            .padding(.vertical, 10)
            .background(Color(.secondarySystemBackground))
            .clipShape(UnevenRoundedRectangle(
                topLeadingRadius: 4, bottomLeadingRadius: 18,
                bottomTrailingRadius: 18, topTrailingRadius: 18
            ))

            Spacer(minLength: 60)
        }
        .onAppear { cursorVisible = false }
    }

    private var markdownContent: AttributedString {
        let lines = content.components(separatedBy: "\n")
        var result = AttributedString()
        for (index, line) in lines.enumerated() {
            result += (try? AttributedString(markdown: line)) ?? AttributedString(line)
            if index < lines.count - 1 { result += AttributedString("\n") }
        }
        return result
    }
}
