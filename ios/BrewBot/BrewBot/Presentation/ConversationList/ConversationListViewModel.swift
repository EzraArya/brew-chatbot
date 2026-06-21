//
//  ConversationListViewModel.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation
import SwiftUI

@Observable
@MainActor
final class ConversationListViewModel {
    var conversations: [Conversation] = []
    private let repository: ChatRepositoryProtocol

    init(repository: ChatRepositoryProtocol) {
        self.repository = repository
    }

    func loadConversations() {
        conversations = repository.loadConversations()
    }

    func deleteConversation(at offSets: IndexSet) {
        offSets.forEach { index in
            repository.deleteConversation(id: conversations[index].id)
        }
        conversations.remove(atOffsets: offSets)
    }

    func newConversation() -> Conversation {
        Conversation()
    }
}