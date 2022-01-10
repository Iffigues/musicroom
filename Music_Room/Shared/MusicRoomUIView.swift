//
//  MusicRoomUIView.swift
//  Music_Room (iOS)
//
//  Created by Arthur Guérin on 30/12/2021.
//

import SwiftUI

//TO DELETE
struct song: Identifiable {
    var id = UUID()
    var name: String
}

struct MusicRoomUIView: View {
    @State var music_paused = false
    
    @State var song_name: String = ""
    
    let playlist = [song(name:"test"), song(name:"1234"), song(name:"hello")
    ]
    
    var body: some View {
        
        VStack {
            SearchUIView(text: $song_name)
            HStack {
                Button {
                    print("Settings button was tapped")
                } label: {
                    Image(systemName: "gear.circle.fill")
                        .resizable(resizingMode: .tile)
                        .foregroundColor(.white)
                        .frame(width: 30.0, height: 30.0)
                }
                Button {
                    print("Settings button was tapped")
                } label: {
                    Image(systemName: "square.and.arrow.up.circle.fill")
                        .resizable(resizingMode: .tile)
                        .foregroundColor(.white)
                        .frame(width: 30, height: 30)
                }
            }
            .padding(.vertical, 10.0)
            List {
                ForEach(playlist) { song in
                HStack {
                    Button {
                        print("Edit button was tapped")
                    } label: {
                        Image(systemName: "arrow.up.arrow.down")
                            .foregroundColor(.blue)
                    }
                    
                    Text("Song name")
                        .foregroundColor(Color.black)
                    Spacer()
                    Button {
                        print("Trash button was tapped")
                    } label: {
                        Image(systemName: "trash.fill")
                    }
                }
                }
            }
            .colorMultiply(/*@START_MENU_TOKEN@*/Color(red: 0.43137254901960786, green: 0.3803921568627451, blue: 0.8862745098039215)/*@END_MENU_TOKEN@*/)
            
            Spacer()
            
            Divider()
            
            VStack(spacing: 8) {
                Text("Song Title")
                    .font(Font.system(.title).bold())
                    .foregroundColor(Color.white)
                Text("Artist Name")
                    .font(.system(.headline))
                    .foregroundColor(Color.white)
            }
            HStack(spacing: 40) {
                Button(action: {
                    print("Rewind")
                }) {
                    ZStack {
                        Circle()
                            .frame(width: 80, height: 80)
                            .accentColor(.pink)
                            .shadow(radius: 10)
                        Image(systemName: "backward.fill")
                            .foregroundColor(.white)
                            .font(.system(.title))
                    }
                }
                
                Button(action: {
                    print("Pause")
                }) {
                    ZStack {
                        Circle()
                            .frame(width: 80, height: 80)
                            .accentColor(.pink)
                            .shadow(radius: 10)
                        if music_paused == true {
                            Image(systemName: "pause.fill")
                                .foregroundColor(.white)
                                .font(.system(.title))
                        }
                        else {                    Image(systemName: "play.fill")
                                .foregroundColor(.white)
                                .font(.system(.title))
                        }
                    }
                }
                
                Button(action: {
                    print("Skip")
                }) {
                    ZStack {
                        Circle()
                            .frame(width: 80, height: 80)
                            .accentColor(.pink)
                            .shadow(radius: 10)
                        Image(systemName: "forward.fill")
                            .foregroundColor(.white)
                            .font(.system(.title))
                    }
                }
            }
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(
            LinearGradient(gradient: Gradient(colors: [.purple, .blue]), startPoint: .top, endPoint: .bottom)
                .edgesIgnoringSafeArea(.all))
    }
}

struct MusicRoomUIView_Previews: PreviewProvider {
    static var previews: some View {
        MusicRoomUIView()
    }
}
