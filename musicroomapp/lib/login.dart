import 'package:flutter/material.dart';
import 'package:musicroomapp/style.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_web_auth/flutter_web_auth.dart';
import 'dart:convert' show jsonDecode;
import 'package:flutter/services.dart';
import 'package:uuid/uuid.dart';

class MyLogin extends StatefulWidget {
  const MyLogin({Key? key}) : super(key: key);

  @override
  MyCustomFormState createState() {
    return MyCustomFormState();
  }
}

class MyCustomFormState extends State<MyLogin> {
  final _formKey = GlobalKey<FormState>();
  String _status = '';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 80.0, vertical: 30.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                'Login or sign in',
                style: Theme.of(context).textTheme.headline1,
              ),
              Row(mainAxisAlignment: MainAxisAlignment.center, children: [
                Text(
                  'please?',
                  style: Theme.of(context).textTheme.headline1,
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(8.0, 0, 0, 10.0),
                  child: ElevatedButton(
                    child: const Text('SIGN IN'),
                    onPressed: () {
                      Navigator.pushReplacementNamed(context, '/signin');
                    },
                    style: ElevatedButton.styleFrom(
                      primary: Colors.yellow,
                    ),
                  ),
                ),
              ]),
              ElevatedButton(
                child: const Text('CONNECT WITH 42'),
                onPressed: () {

                  void authenticate() async {

                    final uuid = Uuid();
                    final state = uuid.v1();
                    final callbackUrlScheme = 'com.school42.musicroom.musicroomapp';
                    final clientId = '3b260c153d6a0269b28a576526290fac14c7ae53b3cd21c0e4b357a110360a39';
                    final clientSecret = '91f87d2a048d00f864c7543d80e08853530d25b173e38569b648c94f62bcd2d4';

                    final url = Uri.https( 'api.intra.42.fr', '/oauth/authorize',
                        {
                          'response_type': 'code',
                          'client_id': clientId,
                          'redirect_uri': '$callbackUrlScheme://playlist',
                          'state': state
                        });
                    try {

                      final result = await FlutterWebAuth.authenticate(url: url.toString(), callbackUrlScheme: callbackUrlScheme);

                      // Extract code from resulting url
                      final code = Uri.parse(result).queryParameters['code'];

                      final urlPost = Uri.https( 'api.intra.42.fr', '/oauth/token');
                      //this code to get an access token
                      final response = await http.post(urlPost, body: {
                        'client_id': clientId,
                        'redirect_uri': '$callbackUrlScheme://playlist',
                        'grant_type': 'authorization_code',
                        'client_secret': clientSecret,
                        'code': code,
                        'state': state,
                      });

                      final accessToken = jsonDecode(response.body)['access_token'] as String;

                      // Send access token to server
                      final urlGet = Uri.http( 'gopiko.fr:9000', '/user/token', {'code': accessToken});
                      final res = await http.get(urlGet);
                      if (accessToken != '') {
                        Navigator.pushReplacementNamed(context, '/playlist');
                      }
                    } on PlatformException catch (e) {
                      setState(() { _status = 'Got error: $e'; });
                    }
                  }
                  authenticate();
                },
                style: ElevatedButton.styleFrom(
                  primary: Colors.yellow,
                ),
              ),
              Text(
                'Login',
                style: Theme.of(context).textTheme.headline5,
              ),
              Form(
                key: _formKey,
                child: Column(
                  children: [
                    TextFormField(
                        decoration: const InputDecoration(hintText: 'Username'),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'enter a username';
                          }
                        }),
                    TextFormField(
                        decoration: const InputDecoration(
                          hintText: 'Password',
                        ),
                        obscureText: true,
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'enter a password';
                          }
                        })
                  ],
                ),
              ),
              const SizedBox(
                height: 24,
              ),
              ElevatedButton(
                child: const Text('ENTER'),
                onPressed: () {
                  // setState(() {
                    if (_formKey.currentState!.validate()) {
                      Navigator.pushReplacementNamed(context, '/playlist');
                    }
                    // });
                },
                style: ElevatedButton.styleFrom(
                  primary: Colors.yellow,
                ),
              ),
              Text(
                _status,
                style: Theme.of(context).textTheme.headline5,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
