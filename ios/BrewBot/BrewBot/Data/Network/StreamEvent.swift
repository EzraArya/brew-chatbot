//
//  StreamEvent.swift
//  BrewBot
//
//  Created by Ezra Arya Wijaya on 21/06/26.
//

enum StreamEvent {
    case textChunk(String)
    case toolCall(type: ToolType, payload: String)
}