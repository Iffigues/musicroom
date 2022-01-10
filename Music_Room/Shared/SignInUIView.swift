//
//  SignInUIView.swift
//  Music_Room
//
//  Created by Arthur Guérin on 29/12/2021.
//

import SwiftUI

struct SignInUIView: View {
    // MARK: - User Info
    @State private var email = ""
    @State private var password = ""
    
    // MARK: - View Details
    var body: some View {
        VStack() {
            Text("Sign in to Music Room")
                .font(.largeTitle).foregroundColor(Color.white)
                .padding([.top, .bottom], 40)
            Image("sample")
                .resizable()
                .frame(width: 250, height: 250)
                .clipShape(Circle())
                .overlay(Circle().stroke(Color.white, lineWidth: 4))
                .shadow(radius: 10)
                .padding(.bottom, 50)
            
            VStack(alignment: .leading, spacing: 15) {
                
                TextField("Email", text: self.$email)
                    .padding()
                    .background(Color.themeText)
                    .cornerRadius(20.0)
                
                SecureField("Password", text: self.$password)
                    .padding()
                    .background(Color.themeText)
                    .cornerRadius(20.0)
            }
            .padding([.leading, .trailing], 27.5)
            
            Button(action: {}) {
                Text("Sign In")
                    .font(.headline)
                    .foregroundColor(.white)
                    .padding()
                    .frame(width: 300, height: 50)
                    .background(Color.green)
                    .cornerRadius(15.0)
            }
            .padding(.top, 50)
            
            Spacer()
            
            HStack(spacing: 0) {
                Text("Don't have an account? ")
                Button(action: {}) {
                    Text("Sign Up")
                        .foregroundColor(.black)
                }
            }
        }
        .background(
            LinearGradient(gradient: Gradient(colors: [.purple, .blue]), startPoint: .top, endPoint: .bottom)
                .edgesIgnoringSafeArea(.all))
    }
}

struct SignInUIView_Previews: PreviewProvider {
    static var previews: some View {
        SignInUIView()
    }
}
