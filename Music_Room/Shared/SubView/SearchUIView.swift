//
//  SearchUIView.swift
//  Music_Room
//
//  Created by Arthur Guérin on 30/12/2021.
//

import SwiftUI

struct SearchUIView: View {
    @Binding var text: String
    @State private var isEditing = false
    
    var body: some View {
        HStack {
            /// Search subview to find rooms
            TextField("Search for...", text: $text)
                .padding(7)
                .padding(.horizontal, 25)
                .background(Color(.systemGray6))
                .cornerRadius(8)
                .overlay(
                    
                    HStack {
                        
                        Image(systemName: "magnifyingglass")
                            .foregroundColor(.gray)
                            .frame(minWidth: 0, maxWidth: .infinity, alignment: .leading)
                            .padding(.leading, 8)
                        
                        if isEditing {
                            
                            Button(action: {
                                self.text = ""
                                
                            }) {
                                
                                Image(systemName: "multiply.circle.fill")
                                    .foregroundColor(.gray)
                                    .padding(.trailing, 8)
                                
                            }
                        }
                    }
                )
                .padding(.horizontal, 10)
                .onTapGesture {
                    self.isEditing = true
                }
            
            if isEditing {
                
                Button(action: {
                    self.isEditing = false
                    self.text = ""
                    
                })
                {
                    
                    Text("Cancel")
                    
                }
                .padding(.trailing, 10)
                .transition(.move(edge: .trailing))
                .animation(/*@START_MENU_TOKEN@*/.default/*@END_MENU_TOKEN@*/, value: 1)
            }
        }
    }
}

struct SearchUIView_Previews: PreviewProvider {
    static var previews: some View {
        SearchUIView(text: .constant(""))
    }
}
