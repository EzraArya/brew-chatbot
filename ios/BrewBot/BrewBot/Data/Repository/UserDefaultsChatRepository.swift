//
//  UserDefaultsChatRepository.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation

final class UserDefaultsChatRepository: ChatRepositoryProtocol {
    private let storageKey = "brewbot.conversations"
    private let defaults: UserDefaults

    init(defaults: UserDefaults = .standard) {
        self.defaults = defaults
    }

    func loadConversations() -> [Conversation] {
        guard let data = defaults.data(forKey: storageKey), let conversations = try? JSONDecoder().decode([Conversation].self, from: data) else {
            return []
        }
        return conversations.sorted { $0.updatedAt > $1.updatedAt }
    }

    func saveConversation(_ conversation: Conversation) {
        var all = loadConversations()

        if let index = all.firstIndex(where: { $0.id == conversation.id }) {
            all[index] = conversation
        } else {
            all.append(conversation)
        }

        if let data = try? JSONEncoder().encode(all) {
            defaults.set(data, forKey: storageKey)
        }
    }

    func deleteConversation(id: UUID) {
        var all = loadConversations()
        all.removeAll { $0.id == id }

        if let data = try? JSONEncoder().encode(all) {
            defaults.set(data, forKey: storageKey)
        }
    }
}