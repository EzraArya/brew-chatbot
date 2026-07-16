//
//  ChatService.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 20/06/26.
//

import Foundation 
import Domain

@MainActor
public final class ChatService {
    public let baseURL: String 

    public init(baseURL: String = "http://localhost:8080") {
        self.baseURL = baseURL
    }

    public func sendMessage(history: [Message], userMessage: String) async throws -> Message {
        guard let url = URL(string: "\(baseURL)/chat") else {
            throw ChatError.serverError
        }

        let requestBody = ChatRequestDTO(
            history: history.map {
                MessageDTO(from: $0)
            },
            userMessage: userMessage
        )

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = try JSONEncoder().encode(requestBody)

        let (data, response) = try await URLSession.shared.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw ChatError.serverError
        }
        let decoded = try JSONDecoder().decode(ChatResponseDTO.self, from: data)
        return Message(role: .model, content: decoded.reply)
    }
}

extension ChatService {
    enum ChatError: Error {
        case invalidURL
        case serverError

        var errorDescription: String? {
            switch self {
            case .invalidURL:
                return "Invalid URL"
            case .serverError:
                return "BrewBot is unavailable. Please try again."
            }
        }
    }
}
