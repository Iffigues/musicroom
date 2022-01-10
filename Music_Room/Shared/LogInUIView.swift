//
//  LogInUIView.swift
//  Music_Room
//
//  Created by Arthur Guérin on 29/12/2021.
//

import SwiftUI
import SafariServices
import BetterSafariView

struct LogInUIView: View {
    // MARK: - User Info
    @State private var email = ""
    @State private var password = ""
    
    
    // 42: - Intra
    @State private var startingWebAuthenticationSession = false
    @State private var showingAlert = false
    @State private var isLoggedin = false
    
    // MARK: - View Details
    var body: some View {
        NavigationView {
            
            VStack() {
                Text("Welcome to Music Room")
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
                
                Button(action: {
                    print("Sign In Button")
                }) {
                    Text("Sign In")
                        .font(.headline)
                        .foregroundColor(.white)
                        .padding()
                        .frame(width: 200, height: 50)
                        .background(Color.green)
                        .cornerRadius(15.0)
                }
                .padding(.top, 10)
                
                Button(action: {
                    print("Sign In with 42 Intra")
                    self.startingWebAuthenticationSession = true
                }) {
                    Text("Sign In with 42 Intra")
                        .font(.headline)
                        .foregroundColor(.white)
                        .padding()
                        .frame(width: 200, height: 50)
                        .background(Color.gray)
                        .cornerRadius(15.0)
                }
                .alert(isPresented: $showingAlert) {
                    Alert(title: Text("🚨 Broken Connection 🚨"), message: Text("Check your Internet connection or allow 42 API connection"), dismissButton: .default(Text("Got it!")))
                }
                .webAuthenticationSession(isPresented: $startingWebAuthenticationSession) {
                    WebAuthenticationSession(
                        url: URL(string: "https://api.intra.42.fr/oauth/authorize?client_id=06c9be585a0111c927a42c6cde1e8f1a595085a5bf2a97d35955cb8a0366e8cd&redirect_uri=MusicRoom42%3A%2F%2Fcallback&response_type=code")!,
                        callbackURLScheme: "redirect_uri"
                    ) { callbackURL, error in
                        if (error != nil) {
                            debugPrint(error as Any)
                            self.showingAlert = true
                            return
                        }
                        let access = authCallToken(code: callbackURL!)
                        if (access == 1) {
                            self.isLoggedin = true
                        }
                        else {
                            debugPrint(error as Any)
                            self.showingAlert = true
                            return
                        }
                    }
                    .prefersEphemeralWebBrowserSession(false)
                }
                .padding(.top, 10)
                
                
                Spacer()
                
                HStack(spacing: 0) {
                    Text("Don't have an account? ")
                    NavigationLink(destination: SignUpUIView()) {
                        Text("Sign Up")
                            .foregroundColor(.black)
                    }
                    .navigationBarHidden(true)
                    .buttonStyle(PlainButtonStyle())
                }
            }
            .background(
                LinearGradient(gradient: Gradient(colors: [.purple, .blue]), startPoint: .top, endPoint: .bottom)
                    .edgesIgnoringSafeArea(.all))
        }
    }
}

extension Color {
    static var themeText: Color {
        return Color(red: 220.0/255.0, green: 230.0/255.0, blue: 230.0/255.0, opacity: 1.0)
    }
}

struct LogInUIView_Previews: PreviewProvider {
    static var previews: some View {
        LogInUIView()
    }
}
