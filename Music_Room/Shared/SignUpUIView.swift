//
//  SignUpUIView.swift
//  Music_Room
//
//  Created by Arthur Guérin on 29/12/2021.
//

import SwiftUI

struct SignUpUIView: View {
    @Environment(\.presentationMode) var presentationMode: Binding<PresentationMode>

    @State private var showingAlert = false

    
    // MARK: - User Info
    @State private var username = ""
    @State private var email = ""
    @State private var password = ""
    
    // MARK: - View Details
    var body: some View {
        VStack() {
            Text("Sign up to Music Room")
                .font(.largeTitle).foregroundColor(Color.white)
                .padding([.bottom], 40)
            Image("New_user")
                .resizable()
                .scaledToFit()
                .frame(width: 100, height: 100)
                .background(Color(red: 0.291, green: 0.306, blue: 0.76))
                .clipShape(Circle())
                .overlay(Circle().stroke(Color.white, lineWidth: 4))
                .shadow(radius: 10)
                .padding(.bottom, 50)
            
            VStack(alignment: .leading, spacing: 15) {
                TextField("Username", text: self.$username)
                    .padding()
                    .background(Color.themeText)
                    .cornerRadius(20.0)
                
                TextField("Location", text: self.$email)
                    .padding()
                    .background(Color.themeText)
                    .cornerRadius(20.0)
                
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
            
            Button(action: {
                print("Sign Up Button")
                showingAlert = true
            }) {
                Text("Sign Up")
                    .font(.headline)
                    .foregroundColor(.white)
                    .padding()
                    .frame(width: 300, height: 50)
                    .background(Color.green)
                    .cornerRadius(15.0)
            }
            .padding(.top, 50)
            .alert("Sign Up error", isPresented: $showingAlert) {
                Button("OK", role: .cancel) { }
            }
            
            Spacer()
            
        }
        .background(
            LinearGradient(gradient: Gradient(colors: [.purple, .blue]), startPoint: .top, endPoint: .bottom)
                .edgesIgnoringSafeArea(.all))
        
        .navigationBarBackButtonHidden(true)
        .navigationBarItems(leading: Button(action : {
            self.presentationMode.wrappedValue.dismiss()
        }){
            Image(systemName: "arrow.left")
            Text("Already have an account?")
        }).foregroundColor(.white)

    }
}

struct SignUpUIView_Previews: PreviewProvider {
    static var previews: some View {
        SignUpUIView()
    }
}
