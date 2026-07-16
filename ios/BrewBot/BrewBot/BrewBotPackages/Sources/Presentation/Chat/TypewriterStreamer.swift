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
    
    private actor StreamState {
        var charBuffer: [Character] = []
        var isNetworkFinished = false
        var networkError: Error?
        var interceptedTool: ToolType? = nil
        var interceptedPayload: String? = nil
        
        func append(_ text: String) { charBuffer.append(contentsOf: text) }
        func setTool(_ type: ToolType, payload: String) {
            interceptedTool = type
            interceptedPayload = payload
        }
        func setFinished(error: Error? = nil) {
            isNetworkFinished = true
            networkError = error
        }
        
        func getStatus(currentIndex: Int) -> (unreadCount: Int, isFinished: Bool, error: Error?, tool: ToolType?, payload: String?) {
            return (charBuffer.count - currentIndex, isNetworkFinished, networkError, interceptedTool, interceptedPayload)
        }
        
        func getChar(at index: Int) -> Character {
            return charBuffer[index]
        }
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
        
        let state = StreamState()
        
        // 1. The Producer: Drains the network as fast as possible
        let producer = Task {
            do {
                for try await event in stream {
                    switch event {
                    case .textChunk(let text):
                        await state.append(text)
                    case .toolCall(let type, let payload):
                        await state.setTool(type, payload: payload)
                        return // Stop draining, we found a tool!
                    }
                }
            } catch {
                await state.setFinished(error: error)
                return
            }
            await state.setFinished()
        }
        
        // 2. The Consumer: Paces the UI updates
        var currentIndex = 0
        while true {
            let status = await state.getStatus(currentIndex: currentIndex)
            
            if let error = status.error {
                producer.cancel()
                throw error
            }
            
            if status.tool != nil {
                break // Abort typing if a tool arrived
            }
            
            if status.unreadCount > 0 {
                let char = await state.getChar(at: currentIndex)
                currentIndex += 1
                
                await onTextAppended(char)
                
                let delay: UInt64 = status.unreadCount > 30 ? 8_000_000 : 15_000_000
                try? await Task.sleep(nanoseconds: delay)
            } else if status.isFinished {
                break // Finished and nothing left to read
            } else {
                // Wait for more data
                try? await Task.sleep(nanoseconds: 10_000_000)
            }
        }
        
        producer.cancel()
        let finalStatus = await state.getStatus(currentIndex: 0)
        return (finalStatus.tool, finalStatus.payload)
    }
}
