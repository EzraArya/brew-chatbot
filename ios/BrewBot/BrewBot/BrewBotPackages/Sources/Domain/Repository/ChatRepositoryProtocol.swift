//
//  ChatRepositoryProtocol.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation

public protocol ChatRepositoryProtocol {
    func loadConversations() -> [Conversation]
    func saveConversation(_ conversation: Conversation)
    func deleteConversation(id: UUID)
}