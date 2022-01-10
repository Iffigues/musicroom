//
//  MusicRoomCreationUIView.swift
//  Music_Room
//
//  Created by Arthur Guérin on 11/01/2022.
//

import SwiftUI

struct MusicRoomCreationUIView: View {
    @State var music_paused = false
    
    @State var song_name: String = ""
    
    let playlist = [song(name:"test"), song(name:"1234"), song(name:"hello")
    ]
    
    var body: some View {
        
        VStack {
            //SearchUIView(text: $song_name)
            HStack {
                Button {
                    print("add button was tapped")
                } label: {
                    Image(systemName: "plus.circle.fill")
                        .resizable(resizingMode: .tile)
                        .foregroundColor(.white)
                        .frame(width: 30.0, height: 30.0)
                }
                
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
                    
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(
            LinearGradient(gradient: Gradient(colors: [.purple, .blue]), startPoint: .top, endPoint: .bottom)
                .edgesIgnoringSafeArea(.all))
    }
}

struct MusicRoomCreationUIView_Previews: PreviewProvider {
    static var previews: some View {
        MusicRoomCreationUIView()
    }
}
