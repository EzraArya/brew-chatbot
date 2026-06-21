//
//  ChatViewModel.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation

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
            let stream = try service.streamMessage(
                history: conversation.messages.dropLast().map { $0 }, 
                userMessage: trimmed
            )

            try await consumeAndAnimateStream(stream)

            let finalMessage = Message(role: .model, content: streamingContent)
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

    private func consumeAndAnimateStream(_ stream: AsyncThrowingStream<String, Error>) async throws {
        var charBuffer: [Character] = []
        var isNetworkFinished = false
        var networkError: Error?
        
        let producer = Task {
            do {
                for try await chunk in stream {
                    charBuffer.append(contentsOf: chunk)
                }
            } catch {
                networkError = error
            }
            isNetworkFinished = true
        }
        
        var currentIndex = 0
        
        while !isNetworkFinished || currentIndex < charBuffer.count {            
            if let error = networkError {
                throw error
            }
            
            let unreadCount = charBuffer.count - currentIndex
            
            if unreadCount > 0 {
                let char = charBuffer[currentIndex]
                streamingContent.append(char)
                currentIndex += 1

                let delayNanoseconds: UInt64 = unreadCount > 30 ? 8_000_000 : 15_000_000
                
                try? await Task.sleep(nanoseconds: delayNanoseconds)
                
            } else {
                try? await Task.sleep(nanoseconds: 10_000_000)
            }
        }
        
        producer.cancel()
    }
}
