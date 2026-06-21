//
//  TypewriterStreamer.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Foundation

enum TypewriterStreamer {
    
    /// Drains a StreamEvent network stream and paces the text output for a UI typewriter effect.
    /// - Parameters:
    ///   - stream: The network stream from ChatService
    ///   - onTextAppended: A closure called on the Main thread every time a character should appear on screen
    /// - Returns: A tuple containing the ToolType and Payload if a tool was intercepted, otherwise nil.
    static func animate(
        stream: AsyncThrowingStream<StreamEvent, Error>,
        onTextAppended: @escaping @MainActor (Character) -> Void
    ) async throws -> (ToolType?, String?) {
        
        var charBuffer: [Character] = []
        var isNetworkFinished = false
        var networkError: Error?
        
        var interceptedTool: ToolType? = nil
        var interceptedPayload: String? = nil
        
        // 1. The Producer: Drains the network as fast as possible
        let producer = Task {
            do {
                for try await event in stream {
                    switch event {
                    case .textChunk(let text):
                        charBuffer.append(contentsOf: text)
                    case .toolCall(let type, let payload):
                        interceptedTool = type
                        interceptedPayload = payload
                        return // Stop draining, we found a tool!
                    }
                }
            } catch {
                networkError = error
            }
            isNetworkFinished = true
        }
        
        // 2. The Consumer: Paces the UI updates
        var currentIndex = 0
        while !isNetworkFinished || currentIndex < charBuffer.count {
            if let error = networkError { throw error }
            
            if interceptedTool != nil { break } // Abort typing if a tool arrived
            
            let unreadCount = charBuffer.count - currentIndex
            if unreadCount > 0 {
                let char = charBuffer[currentIndex]
                currentIndex += 1
                
                // Jump to the Main Thread to update the UI
                await MainActor.run {
                    onTextAppended(char)
                }
                
                let delay: UInt64 = unreadCount > 30 ? 8_000_000 : 15_000_000
                try? await Task.sleep(nanoseconds: delay)
            } else {
                try? await Task.sleep(nanoseconds: 10_000_000)
            }
        }
        
        producer.cancel()
        return (interceptedTool, interceptedPayload)
    }
}