//
//  SpotifyConnectionUIView.swift
//  Music_Room
//
//  Created by Arthur Guérin on 29/12/2021.
//

import SwiftUI

struct SpotifyConnectionUIView: View {
    var body: some View {
            VStack() {
                Text("Connect to Spotify")
                    .font(.largeTitle).foregroundColor(Color.white)
                    .padding(.vertical, 90)
                Image("spotify")
                    .resizable()
                    .scaledToFit()
                    .frame(width: 100, height: 100)
                    .background(Color(red: 0.291, green: 0.306, blue: 0.76))
                    .clipShape(Circle())
                    .overlay(Circle().stroke(Color.white, lineWidth: 4))
                    .shadow(radius: 10)
                    .padding(.bottom, 50)
                
                Button(action: {
                    print("Spotify")
                }) {
                    Text("Connect your account")
                        .font(.headline)
                        .foregroundColor(.white)
                        .padding()
                        .frame(width: 300, height: 50)
                        .background(Color.green)
                        .cornerRadius(15.0)
                }
                .padding(.top, 30)
                .padding([.leading, .trailing], 27.5)
                
                Spacer()
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(
            LinearGradient(gradient: Gradient(colors: [.purple, .blue]), startPoint: .top, endPoint: .bottom)
                .edgesIgnoringSafeArea(.all))
    }
    
}

struct SpotifyConnectionUIView_Previews: PreviewProvider {
    static var previews: some View {
        SpotifyConnectionUIView()
    }
}
