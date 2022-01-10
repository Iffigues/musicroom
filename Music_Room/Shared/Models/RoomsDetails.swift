//
//  RoomsDetails.swift
//  Music_Room
//
//  Created by Arthur Guérin on 30/12/2021.
//

import Foundation

/// Room  Structure
//struct Room: Identifiable {
//    var id = UUID()
//    var name: String = ""
//}

struct Room: Codable, Identifiable {
    var id = UUID()
    var name: String
    var creator_id: Int
    var song: String
    var playlist: String
    var current_position: Int
    var playlist_type: Int
}

// TO DELETE - SAMPLE PURPOSE ONLY
/// Get the data  from txt file and move them to a Room Array
func loadR(filename: String ) -> [Room]  {
    var Rooms: [Room] = []
    let url = Bundle.main.url(forResource: filename, withExtension: "txt")!
    let data = try! Data(contentsOf: url)
    let string = String(data: data, encoding: .utf8)!
    let room_name = string.components(separatedBy: "\n")
    var a = 0

    room_name.forEach{ elem in
        Rooms.append(Room(name: elem, creator_id: 1, song: "", playlist: "", current_position: 0, playlist_type: 1))
        a += 1
    }
    return Rooms
}
