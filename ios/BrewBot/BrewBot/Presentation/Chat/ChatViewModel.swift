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
    var messages: [Message] = []
    var isLoading: Bool = false
    var errorMessage: String? = nil
    
    private let service: ChatService
    
    init(service: ChatService) {
        self.service = service
    }

    func sendMessage(_ message: String) async {
        let trimmed = message.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty, !isLoading else { return }

        let userMessage = Message(role: .user, content: trimmed)
        messages.append(userMessage)
        
        isLoading = true
        errorMessage = nil

        do {
            let reply = try await service.sendMessage(
                history: messages.dropLast().map {$0}, 
                userMessage: trimmed
            )
            messages.append(reply)
        } catch {
            messages.removeLast()
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }
}
