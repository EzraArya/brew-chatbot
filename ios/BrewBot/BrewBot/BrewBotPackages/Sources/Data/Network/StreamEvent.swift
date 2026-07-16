//
//  StreamEvent.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

import Domain

public enum StreamEvent: Sendable {
    case textChunk(String)
    case toolCall(type: ToolType, payload: String)
}
