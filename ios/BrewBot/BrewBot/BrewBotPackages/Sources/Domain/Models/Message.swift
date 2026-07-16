//
//  Message.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 20/06/26.
//

import Foundation

struct Message: Identifiable, Equatable, Sendable, Codable, Hashable {
    let id: UUID 
    let role: Role 
    let content: String
    let timestamp: Date
    var toolType: ToolType?
    var toolPayload: String?

    init(role: Role, content: String, toolType: ToolType? = nil, toolPayload: String? = nil) {
        self.id = UUID()
        self.role = role
        self.content = content
        self.timestamp = Date()
        self.toolType = toolType
        self.toolPayload = toolPayload
    }
}

extension Message {
    enum Role: String, Codable, Sendable {
        case user = "user"
        case model = "model"
    }
}