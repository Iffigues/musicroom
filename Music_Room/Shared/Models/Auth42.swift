//
//  Auth42.swift
//  Music_Room
//
//  Created by Arthur Guérin on 10/01/2022.
//

import Foundation
import UIKit

struct Login: Encodable {
    let grant_type = "client_credentials"
    let client_id: String
    let client_secret: String
    let code: String
}

struct Answer: Codable {
    var access_token: String
    var created_at: Int
    var expires_in: Int
    var scope: String
    var token_type: String
}

var requestAnswer = Answer(
    access_token: "",
    created_at: 0,
    expires_in: 0,
    scope: "",
    token_type: "")

func isValidToken(expires_in: Int) -> Bool {
    if (expires_in <= 0) {
        debugPrint("Token expired")
        return(false)
    }
    else {
        debugPrint("Token valid")
        return(true)
    }
}

func authCallToken(code: URL) -> Int {
    var dedicatedCode: String
    
    if (code == nil) {
        debugPrint("Error 0: Connection issue")
        return(0)
    }
    
    dedicatedCode = code.absoluteString
    dedicatedCode = dedicatedCode.replacingOccurrences(
        of: "MusicRoom42://callback?code=",
        with: ""
    )
    
    if (dedicatedCode == "MusicRoom42://callback?error=access_denied&error_description=The+resource+owner+or+authorization+server+denied+the+request.")
    {
        debugPrint("Error 1: API access refused")
        return(0)
    }
    
    if (isValidToken(expires_in: requestAnswer.expires_in) == false) {
        let login = Login(
            client_id: "06c9be585a0111c927a42c6cde1e8f1a595085a5bf2a97d35955cb8a0366e8cd",
            client_secret: "a53a03d2e47af307639d344326f714d557fb234de469730418b3fa206d7a6020",
            code: dedicatedCode
        )
        
        var request = URLRequest(url: URL(string: "https://api.intra.42.fr/oauth/token")!)
        request.httpMethod = "post"
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        do {
            let jsonData = try JSONEncoder().encode(login)
            request.httpBody = jsonData
        }catch let jsonErr{
            print(jsonErr)
        }
        URLSession.shared.dataTask(with: request) { data, response, error in
            if let data = data {
                print(String(data: data, encoding: .utf8)!)
                if let response = try?
                    JSONDecoder().decode(Answer.self, from: data) {
                    DispatchQueue.main.async {
                        requestAnswer = response
                        print("AAAAA")
                        print(requestAnswer)
                    }
                    return
                }
                else {
                    debugPrint("Error 2: API Connection issue")
                }
            }
            
        }.resume()
    }
    return(1)
}
