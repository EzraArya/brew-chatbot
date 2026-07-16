//
//  TypewriterStreamer.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation
import Domain
import Data

public enum TypewriterStreamer {

    /// Everything the consumer needs about a produced chunk, including a
    /// cumulative send-count so pacing can react to true backlog depth —
    /// without either side owning shared mutable state to compute it.
    private enum Signal: Sendable {
        case chunk(characters: [Character], totalSentSoFar: Int)
        case toolIntercepted(ToolType, String?)
    }

    /// Drains a StreamEvent network stream and paces the text output for a UI typewriter effect.
    /// - Parameters:
    ///   - stream: The network stream from ChatService
    ///   - onTextAppended: A closure called on the Main thread every time a character should appear on screen
    /// - Returns: A tuple containing the ToolType and Payload if a tool was intercepted, otherwise nil.
    public static func animate(
        stream: AsyncThrowingStream<StreamEvent, Error>,
        onTextAppended: @escaping @MainActor @Sendable (Character) -> Void
    ) async throws -> (ToolType?, String?) {

        let (signals, continuation) = AsyncThrowingStream<Signal, Error>.makeStream()

        // Producer: owns `totalSent` exclusively. Single writer -> no
        // synchronization needed, no actor required for this counter.
        let producer = Task {
            var totalSent = 0
            do {
                for try await event in stream {
                    switch event {
                    case .textChunk(let text):
                        totalSent += text.count
                        continuation.yield(.chunk(characters: Array(text), totalSentSoFar: totalSent))
                    case .toolCall(let type, let payload):
                        continuation.yield(.toolIntercepted(type, payload))
                        continuation.finish()
                        return // Stop draining, we found a tool!
                    }
                }
                continuation.finish()
            } catch {
                continuation.finish(throwing: error)
            }
        }
        defer { producer.cancel() }

        // Consumer: owns `totalConsumed` exclusively. Single writer -> safe
        // by construction. `for try await` suspends properly between yields;
        // no manual polling, no shared buffer, no actor hops per character.
        var totalConsumed = 0
        for try await signal in signals {
            switch signal {
            case .chunk(let characters, let totalSentSoFar):
                for char in characters {
                    await onTextAppended(char)
                    totalConsumed += 1
                    let unreadCount = totalSentSoFar - totalConsumed
                    let delay: UInt64 = unreadCount > 30 ? 8_000_000 : 15_000_000
                    try? await Task.sleep(nanoseconds: delay)
                }
            case .toolIntercepted(let tool, let payload):
                return (tool, payload) // Abort typing if a tool arrived
            }
        }

        return (nil, nil)
    }
}
