//
//  Message.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 20/06/26.
//

import Foundation

public struct Message: Identifiable, Equatable, Sendable, Codable, Hashable {
    public let id: UUID 
    public let role: Role 
    public let content: String
    public let timestamp: Date
    public var toolType: ToolType?
    public var toolPayload: String?

    public init(id: UUID = UUID(), role: Role, content: String, timestamp: Date = Date(), toolType: ToolType? = nil, toolPayload: String? = nil) {
        self.id = id
        self.role = role
        self.content = content
        self.timestamp = timestamp
        self.toolType = toolType
        self.toolPayload = toolPayload
    }
}

public extension Message {
    enum Role: String, Codable, Sendable {
        case user = "user"
        case model = "model"
    }
}