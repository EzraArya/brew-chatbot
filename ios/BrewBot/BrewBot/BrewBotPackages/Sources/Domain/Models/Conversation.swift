//
//  Conversation.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation

public struct Conversation: Identifiable, Equatable, Sendable, Codable, Hashable {
    public let id: UUID
    public var title: String
    public var messages: [Message]
    public let createdAt: Date
    public var updatedAt: Date

    public init(id: UUID = UUID(), title: String = "New Conversation", messages: [Message] = [], createdAt: Date = Date(), updatedAt: Date = Date()) {
        self.id = id
        self.title = title
        self.messages = messages
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}