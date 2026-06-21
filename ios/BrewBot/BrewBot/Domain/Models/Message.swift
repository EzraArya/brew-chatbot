//
//  Message.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 20/06/26.
//

import Foundation

struct Message: Identifiable, Equatable, Sendable, Codable {
    let id: UUID 
    let role: Role 
    let content: String
    let timestamp: Date

    init(role: Role, content: String) {
        self.id = UUID()
        self.role = role
        self.content = content
        self.timestamp = Date()
    }
}

extension Message {
    enum Role: String, Codable, Sendable {
        case user = "user"
        case model = "model"
    }
}