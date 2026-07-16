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
    ) -> AsyncThrowingStream<StreamEvent, Error> {
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

                            if chunk.hasPrefix("["), let bracketEnd = chunk.firstIndex(of: "]") {
                                let toolNameRaw = String(chunk[chunk.index(after: chunk.startIndex)..<bracketEnd])
                                if let toolType = ToolType(rawValue: toolNameRaw) {
                                    let payloadRaw = String(chunk[chunk.index(bracketEnd, offsetBy: 2)...])
                                    continuation.yield(.toolCall(type: toolType, payload: payloadRaw))
                                    continue
                                }
                            }

                            let cleanChunk = chunk.replacingOccurrences(of: "\\n", with: "\n")
                            continuation.yield(.textChunk(cleanChunk))
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