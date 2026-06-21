//
//  ChatView.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import SwiftUI

struct ChatView: View {
    private let conversation: Conversation
    private let repository: ChatRepositoryProtocol
    
    init(conversation: Conversation, repository: ChatRepositoryProtocol) {
        self.conversation = conversation
        self.repository = repository
    }
    
    @State private var viewModel: ChatViewModel?
    @State private var inputText = ""
    
    var body: some View {
        VStack(spacing: 0) {
            if let viewModel {
                ScrollViewReader { proxy in
                    ScrollView {
                        LazyVStack(spacing: 8) {
                            ForEach(viewModel.messages) { message in
                                if let type = message.toolType, let payload = message.toolPayload {
                                    
                                    WidgetRenderer(model: WidgetFactory.decode(type: type, payload: payload))
                                        .id(message.id)
                                        .padding(.vertical, 4)
                                        
                                } else {
                                    MessageBubble(message: message)
                                        .id(message.id)
                                }
                            }

                            if viewModel.isStreaming {
                                StreamingBubble(content: viewModel.streamingContent)
                                    .id("streaming")
                            }

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
                    .onAppear {
                        scrollToBottom(proxy: proxy)
                    }
                    .onChange(of: viewModel.isLoading) {
                        scrollToBottom(proxy: proxy)
                    }
                    .onChange(of: viewModel.streamingContent.count) {
                        scrollToBottom(proxy: proxy)
                    }
                }

                Divider()

                InputBar(
                    text: $inputText,
                    isLoading: viewModel.isLoading || viewModel.isStreaming
                ) {
                    sendMessage()
                }
            } else {
                ProgressView()
            }
        }
        .navigationTitle(viewModel?.conversation.title ?? "BrewBot 🍺")
        .navigationBarTitleDisplayMode(.inline)
        .alert("Something went wrong",
               isPresented: Binding(
                get: { viewModel?.errorMessage != nil },
                set: { if !$0 { viewModel?.errorMessage = nil } }
               )) {
            Button("OK", role: .cancel) { }
        } message: {
            Text(viewModel?.errorMessage ?? "")
        }
        .task {
            if viewModel == nil {
                viewModel = ChatViewModel(
                    conversation: conversation,
                    service: ChatService(),
                    repository: repository
                )
            }
        }
    }

    // MARK: - Actions

    private func sendMessage() {
        guard let viewModel else { return }
        let text = inputText
        inputText = ""
        Task {
            await viewModel.streamMessage(text)
        }
    }

    private func scrollToBottom(proxy: ScrollViewProxy) {
        guard let viewModel else { return }
        
        if viewModel.isStreaming {
            proxy.scrollTo("streaming", anchor: .bottom)
        } else {
            withAnimation(.easeOut(duration: 0.3)) {
                if viewModel.isLoading {
                    proxy.scrollTo("typing", anchor: .bottom)
                } else if let last = viewModel.messages.last {
                    proxy.scrollTo(last.id, anchor: .bottom)
                }
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
