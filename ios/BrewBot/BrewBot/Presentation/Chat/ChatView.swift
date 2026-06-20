//
//  ChatView.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import SwiftUI

struct ChatView: View {
    @State private var viewModel = ChatViewModel(service: ChatService())
    @State private var inputText = ""

    var body: some View {
        NavigationStack {
            VStack(spacing: 0) {
                // Message list
                ScrollViewReader { proxy in
                    ScrollView {
                        LazyVStack(spacing: 8) {
                            ForEach(viewModel.messages) { message in
                                MessageBubble(message: message)
                                    .id(message.id)
                            }

                            // Typing indicator
                            if viewModel.isLoading {
                                TypingIndicator()
                                    .id("typing")
                            }
                        }
                        .padding()
                    }
                    .onChange(of: viewModel.messages.count) {
                        scrollToBottom(proxy: proxy)
                    }
                    .onChange(of: viewModel.isLoading) {
                        scrollToBottom(proxy: proxy)
                    }
                }

                Divider()

                // Input bar
                InputBar(
                    text: $inputText,
                    isLoading: viewModel.isLoading
                ) {
                    sendMessage()
                }
            }
            .navigationTitle("BrewBot 🍺")
            .navigationBarTitleDisplayMode(.inline)
            .alert("Something went wrong",
                   isPresented: Binding(
                    get: { viewModel.errorMessage != nil },
                    set: { if !$0 { viewModel.errorMessage = nil } }
                   )) {
                Button("OK", role: .cancel) { }
            } message: {
                Text(viewModel.errorMessage ?? "")
            }
        }
    }

    // MARK: - Actions

    private func sendMessage() {
        let text = inputText
        inputText = ""   // clear immediately for good UX
        Task {
            await viewModel.sendMessage(text)
        }
    }

    private func scrollToBottom(proxy: ScrollViewProxy) {
        withAnimation(.easeOut(duration: 0.3)) {
            if viewModel.isLoading {
                proxy.scrollTo("typing", anchor: .bottom)
            } else if let last = viewModel.messages.last {
                proxy.scrollTo(last.id, anchor: .bottom)
            }
        }
    }
}

// MARK: - Typing Indicator

private struct TypingIndicator: View {
    @State private var opacity = 0.3

    var body: some View {
        HStack {
            HStack(spacing: 4) {
                ForEach(0..<3, id: \.self) { index in
                    Circle()
                        .frame(width: 8, height: 8)
                        .opacity(opacity)
                        .animation(
                            .easeInOut(duration: 0.6)
                            .repeatForever()
                            .delay(Double(index) * 0.2),
                            value: opacity
                        )
                }
            }
            .padding(.horizontal, 14)
            .padding(.vertical, 10)
            .background(Color(.secondarySystemBackground), in: Capsule())

            Spacer(minLength: 60)
        }
        .onAppear { opacity = 1.0 }
    }
}