//
//  SettingsUIView.swift
//  Music_Room (iOS)
//
//  Created by Arthur Guérin on 30/12/2021.
//

import SwiftUI

struct SettingsUIView: View {
    var body: some View {
        VStack(alignment: .leading, spacing: 15) {
            
            Text("Your profile")
                .padding()
                .background(Color.themeText)
                .cornerRadius(20.0)
            
            Text("Your rooms")
                .padding()
                .background(Color.themeText)
                .cornerRadius(20.0)
        }
        .padding([.leading, .trailing], 27.5)
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(
            LinearGradient(gradient: Gradient(colors: [.purple, .blue]), startPoint: .top, endPoint: .bottom)
                .edgesIgnoringSafeArea(.all))
    }
    
}

struct SettingsUIView_Previews: PreviewProvider {
    static var previews: some View {
        SettingsUIView()
    }
}
