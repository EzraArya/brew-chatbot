//
//  ConversationListView.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import SwiftUI

struct ConversationListView: View {

    @State private var viewModel = ConversationListViewModel(
        repository: UserDefaultsChatRepository()
    )
    @State private var newConversation: Conversation? = nil

    var body: some View {
        NavigationStack {
            Group {
                if viewModel.conversations.isEmpty {
                    emptyState
                } else {
                    list
                }
            }
            .navigationTitle("BrewBot 🍺")
            .toolbar {
                ToolbarItem(placement: .primaryAction) {
                    Button {
                        newConversation = viewModel.newConversation()
                    } label: {
                        Image(systemName: "square.and.pencil")
                    }
                }
            }
            .navigationDestination(item: $newConversation) { conversation in
                ChatView(
                    conversation: conversation,
                    repository: viewModel.repository
                )
            }
            .onAppear {
                viewModel.loadConversations()
            }
        }
    }

    // MARK: - Subviews

    private var list: some View {
        List {
            ForEach(viewModel.conversations) { conversation in
                NavigationLink {
                    ChatView(
                        conversation: conversation,
                        repository: viewModel.repository
                    )
                } label: {
                    row(for: conversation)
                }
            }
            .onDelete { offsets in
                viewModel.deleteConversation(at: offsets)
            }
        }
    }

    private var emptyState: some View {
        VStack(spacing: 16) {
            Image(systemName: "bubble.left.and.bubble.right")
                .font(.system(size: 60))
                .foregroundStyle(.secondary)

            Text("No conversations yet")
                .font(.title3)
                .fontWeight(.semibold)

            Text("Tap the pencil to start chatting with BrewBot")
                .font(.subheadline)
                .foregroundStyle(.secondary)
                .multilineTextAlignment(.center)
        }
        .padding()
    }

    private func row(for conversation: Conversation) -> some View {
        VStack(alignment: .leading, spacing: 4) {
            Text(conversation.title)
                .font(.headline)
                .lineLimit(1)

            HStack {
                Text("\(conversation.messages.count) messages")
                Spacer()
                Text(conversation.updatedAt.formatted(.relative(presentation: .named)))
            }
            .font(.caption)
            .foregroundStyle(.secondary)
        }
        .padding(.vertical, 4)
    }
}