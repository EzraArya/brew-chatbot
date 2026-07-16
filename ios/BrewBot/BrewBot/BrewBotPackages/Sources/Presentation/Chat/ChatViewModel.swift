//
//  ChatViewModel.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation
import Domain
import Data

@Observable
@MainActor
final class ChatViewModel {

    // MARK: - State

    private(set) var conversation: Conversation
    var isLoading: Bool = false
    var errorMessage: String? = nil

    var isStreaming: Bool = false
    var streamingContent: String = ""

    /// Convenience — Views read messages from here directly
    var messages: [Message] { conversation.messages }

    // MARK: - Dependencies

    private let service: ChatService
    private let repository: ChatRepositoryProtocol

    init(
        conversation: Conversation,
        service: ChatService,
        repository: ChatRepositoryProtocol
    ) {
        self.conversation = conversation
        self.service = service
        self.repository = repository
    }

    // MARK: - Actions

    func sendMessage(_ text: String) async {
        let trimmed = text.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty, !isLoading else { return }

        let userMessage = Message(role: .user, content: trimmed)
        conversation.messages.append(userMessage)

        generateTitleIfNeeded(from: trimmed)

        isLoading = true
        errorMessage = nil

        do {
            let reply = try await service.sendMessage(
                history: conversation.messages.dropLast().map { $0 },
                userMessage: trimmed
            )

            conversation.messages.append(reply)

        } catch {
            conversation.messages.removeLast()
            errorMessage = "Couldn't reach BrewBot. Please try again."
        }

        save()

        isLoading = false
    }

    func streamMessage(_ text: String) async {
        let trimmed = text.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty, !isLoading, !isStreaming else { return }

        let userMessage = Message(role: .user, content: trimmed)
        conversation.messages.append(userMessage)

        generateTitleIfNeeded(from: trimmed)

        isStreaming = true
        streamingContent = ""
        errorMessage = nil

        do {
            let stream = service.streamMessage(
                history: conversation.messages.dropLast().map { $0 }, 
                userMessage: trimmed
            )

            let (toolType, toolPayload) = try await TypewriterStreamer.animate(stream: stream) { char in
                self.streamingContent.append(char)
            }
            
            let finalMessage = Message(
                role: .model, 
                content: streamingContent,
                toolType: toolType,
                toolPayload: toolPayload
            )
            conversation.messages.append(finalMessage)
        } catch {
            conversation.messages.removeLast()
            errorMessage = error.localizedDescription
        }

        isStreaming = false
        streamingContent = ""
        save()
    }

    // MARK: - Private Helpers

    private func save() {
        conversation.updatedAt = Date()
        repository.saveConversation(conversation)
    }

    private func generateTitleIfNeeded(from text: String) {
        guard conversation.messages.count == 1 else { return }
        let words = text.split(separator: " ").prefix(5).joined(separator: " ")
        conversation.title = words.isEmpty ? "New Conversation" : words
    }
}
