//
//  DTO.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 20/06/26.
//

import Foundation
import Domain

struct MessageDTO: Codable, Sendable {
    let role: String
    let content: String
}

struct ChatRequestDTO: Codable, Sendable {
    let history: [MessageDTO]
    let userMessage: String
}

struct ChatResponseDTO: Codable, Sendable {
    let reply: String
}

extension MessageDTO {
    init(from message: Message) {
        self.role = message.role.rawValue
        self.content = message.content
    }
}

extension Message {
    init(from dto: MessageDTO) {
        self.init(
            role: dto.role == "user" ? .user : .model,
            content: dto.content
        )
    }
}