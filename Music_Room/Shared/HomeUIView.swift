//
//  HomeUIView.swift
//  Music_Room
//
//  Created by Arthur Guérin on 29/12/2021.
//

import SwiftUI

struct HomeUIView: View {
    @State var room_name: String = ""
    @State var rooms = [Room]()
    
    init() {
        //Use this if NavigationBarTitle is with Large Font
        UINavigationBar.appearance().largeTitleTextAttributes = [.foregroundColor: UIColor.white]

        //Use this if NavigationBarTitle is with displayMode = .inline
        UINavigationBar.appearance().titleTextAttributes = [.foregroundColor: UIColor.white]
    }
    
    var body: some View {
        /// Load from API and display the list of  rooms
        NavigationView {
            ZStack {
                LinearGradient(gradient: Gradient(colors: [.purple, .blue]), startPoint: .top, endPoint: .bottom)
                    .edgesIgnoringSafeArea(.all)
                
                VStack {
                    
                    SearchUIView(text: $room_name)
                    
                    Spacer()
                    
                    List(rooms.filter({ room_name.isEmpty ? true : $0.name.contains(room_name) })) { item in
                        NavigationLink(destination: SignUpUIView()) {
                            Text(item.name)
                        }
                    }
                    .ignoresSafeArea()
                    Button(action: {
                        print("New Room")
                    }) {
                        Text("Create a new Room")
                            .font(.headline)
                            .foregroundColor(.white)
                            .padding()
                            .frame(width: 200, height: 50)
                            .background(Color.pink)
                            .cornerRadius(15.0)
                    }

                }
                
                .navigationTitle("Available Rooms")
            }
        }
        //.navigationBarTitle("Music Room - \(room_name)")
        //.navigationBarBackButtonHidden(true)
        //.navigationBarHidden(true)
        .onAppear(perform: {
            //rooms = loadR(filename: "Rooms_Sample")
            print("hello")
            loadData()
            print(rooms)
            print(rooms.count)
        })
    }
    
    func loadData() {
            guard let url = URL(string: "http://gopiko.fr:9000/room/2") else {
                print("Invalid URL")
                return
            }
            var request = URLRequest(url: url)
            request.httpMethod = "GET"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            print(request)
            URLSession.shared.dataTask(with: request) { data, response, error in
                if let data = data {
                    print(String(data: data, encoding: .utf8)!)
                    if let response = try?
                        JSONDecoder().decode([Room].self, from: data) {
                        DispatchQueue.main.async {
                            self.rooms = response
                            print("AAAAA")
                            print(rooms)
                        }
                        return
                    }
                    else {
                        print("ERROR")
                    }
                }

            }.resume()
        }
}

struct HomeUIView_Previews: PreviewProvider {
    static var previews: some View {
        HomeUIView()
    }
}
