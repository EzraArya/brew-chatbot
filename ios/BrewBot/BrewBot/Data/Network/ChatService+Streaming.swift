//
//  ChatService+Streaming.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation

extension ChatService {
    func streamMessage(
        history: [Message],
        userMessage: String
    ) -> AsyncThrowingStream<String, Error> {
        return AsyncThrowingStream { continuation in 
            Task {
                do {
                    let requestDTO = ChatRequestDTO(
                        history: history.map { MessageDTO(from: $0) }, 
                        userMessage: userMessage
                    )

                    guard let url = URL(string: "\(baseURL)/chat/stream") else {
                        throw ChatError.invalidURL
                    }
                    var request = URLRequest(url: url)
                    request.httpMethod = "POST"
                    request.setValue("application/json", forHTTPHeaderField: "Content-Type")
                    request.httpBody = try JSONEncoder().encode(requestDTO)

                    let (bytes, response) = try await URLSession.shared.bytes(for: request)

                    guard let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 else {
                        throw ChatError.serverError
                    }

                    for try await line in bytes.lines {
                        guard !line.isEmpty else { continue }

                        if line.hasPrefix("data: ") {
                            let chunk = String(line.dropFirst(6))

                            if chunk == "[DONE]" {
                                continuation.finish()
                                return
                            }

                            if chunk == "[ERROR]" {
                                continuation.finish(throwing: ChatError.serverError)
                                return
                            }

                            let cleanChunk = chunk.replacingOccurrences(of: "\\n", with: "\n")
                            continuation.yield(cleanChunk)
                        }
                    }

                    continuation.finish()
                } catch {
                    continuation.finish(throwing: error)
                }
            }
        }
    }
}