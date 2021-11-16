import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'dart:math';

import 'package:http/http.dart' as http;

import 'package:flutter/material.dart';
import 'package:musicroomapp/style.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:audioplayers/audioplayers.dart';
import 'package:file_picker/file_picker.dart';
import 'package:spotify_sdk/spotify_sdk.dart';

String getRandString(int len) {
  var random = Random.secure();
  var values = List<int>.generate(len, (i) =>  random.nextInt(255));
  return base64UrlEncode(values);
}

Future<identified> apiCall() async {
  //URL WITH WEBVIEW

  var client_id_0 = '4dee364910944d2a9adcc80b54b28942';
  var redirect_uri_0 = 'http://localhost:8888/callback';

  var state_0 = getRandString(16);
  var scope_0 = 'user-read-private user-read-email streaming user-read-currently-playing';

  var other_auth = 'https://accounts.spotify.com/authorize';
  other_auth += '?response_type=token';
  other_auth += '&client_id=' + Uri.encodeQueryComponent(client_id_0);
  other_auth += '&scope=' + Uri.encodeQueryComponent(scope_0);
  other_auth += '&redirect_uri=' + Uri.encodeQueryComponent(redirect_uri_0);
  other_auth += '&state=' + Uri.encodeQueryComponent(state_0);
  debugPrint(other_auth);
  sleep(Duration(seconds:5));
  http.Request req = http.Request("Get", Uri.parse(other_auth))..followRedirects = false;
  http.Client baseClient = http.Client();
  http.StreamedResponse response_0 = await baseClient.send(req);
  Uri redirectUri = Uri.parse(response_0.headers['location'].toString());
  final response_1 = await http.post(redirectUri);
  debugPrint(response_1.body.toString());
  debugPrint(redirectUri.toString());

  var client_id = '4dee364910944d2a9adcc80b54b28942';
  var client_secret = '739d19ae0f744531b5a5469ff762afcf';
  var bytes = utf8.encode(client_id + ':' + client_secret);
  var base64Str = 'Basic ' + base64.encode(bytes);
  debugPrint(base64Str);
  var data = {
    'grant_type': 'client_credentials'
  };
  final response = await http.post(
    Uri.parse("https://accounts.spotify.com/api/token"),
    // Send authorization headers to the backend.
    headers: {
      HttpHeaders.authorizationHeader: base64Str,
      HttpHeaders.contentTypeHeader: 'application/x-www-form-urlencoded',
      },
      body: data
  );
  sleep(Duration(seconds:5));
  final responseJson = jsonDecode(response.body);
  debugPrint(response.body);
  return identified.fromJson(responseJson);
}

class identified {
  final String token;
  final int expires_in;

  identified({
    required this.token,
    required this.expires_in,
  });

  factory identified.fromJson(Map<String, dynamic> json) {
    return identified(
      token: json['access_token'],
      expires_in: json['expires_in'],
    );
  }

}

var SpotifyIcon = Icon(
  FontAwesomeIcons.spotify,
  color: Colors.black,
);

void foo() async {
  final auth = await apiCall();
  debugPrint(auth.token);
}

// in case we need to sync something with spotify first
var SpotifyButton = IconButton(icon: SpotifyIcon, onPressed: () {
  debugPrint('movieTitle:');
  foo();
});

class MyPlaylist extends StatelessWidget {
  const MyPlaylist({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        home: AudioPlayerScreen(),
        theme: ThemeData(primaryColor: Colors.black));
  }
}

class AudioPlayerScreen extends StatefulWidget {
  @override
  _AudioPlayerScreenState createState() => _AudioPlayerScreenState();
}

class _AudioPlayerScreenState extends State<AudioPlayerScreen> {
  //ROOM LIST
  TextEditingController editingController = TextEditingController();
  final duplicateItems = List<String>.generate(100, (i) => "Element $i");
  var items = <String>[];

  //AUDIO PLAYER
  AudioPlayer _audioPlayer = AudioPlayer();
  bool _isPlaying = false;
  String currentTime = "00:00";
  String completeTime = "00:00";

  @override
  void initState() {
    items.addAll(duplicateItems);
    super.initState();
    _audioPlayer.onAudioPositionChanged.listen((Duration duration) {
      setState(() {
        currentTime = duration.toString().split(".")[0];
      });
    });
    _audioPlayer.onDurationChanged.listen((Duration duration) {
      setState(() {
        completeTime = duration.toString().split(".")[0];
      });
    });
  }

  void filterSearchResults(String query) {
    List<String> dummySearchList = <String>[];
    dummySearchList.addAll(duplicateItems);
    if (query.isNotEmpty) {
      List<String> dummyListData = <String>[];
      dummySearchList.forEach((item) {
        if (item.contains(query)) {
          dummyListData.add(item);
        }
      });
      setState(() {
        items.clear();
        items.addAll(dummyListData);
      });
      return;
    } else {
      setState(() {
        items.clear();
        items.addAll(duplicateItems);
      });
    }
  }

  var playerbar = AppBar(
    title: Container(
      alignment: Alignment.center,
      child: Row(children: <Widget>[
        Image(
          image:
              NetworkImage("https://www.arthurguerin.com/assets/musicroom.png"),
          width: 30,
          height: 30,
          fit: BoxFit.cover,
        ),
        Padding(
          padding: EdgeInsets.all(16.0),
          child: Text(
            'Music Room',
            style: TextStyle(
              color: Colors.black,
            ),
            textDirection: TextDirection.ltr,
          ),
        ),
      ]),
    ),
    backgroundColor: Colors.yellow,
    actions: <Widget>[],
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.yellow,
      resizeToAvoidBottomInset: false,
      appBar: playerbar,
      body: Container(
        width: MediaQuery.of(context).size.width,
        child: Column(
          children: <Widget>[
            Container(
              decoration: new BoxDecoration(color: Colors.white),
              padding: const EdgeInsets.all(10.0),
              child: Column(children: <Widget>[
                Text(
                  "AVAILABLE ROOM",
                  style: TextStyle(
                      fontSize: 20,
                      color: Colors.black,
                      fontWeight: FontWeight.bold),
                ),
                Container(
                  child: Column(
                    children: <Widget>[
                      Padding(
                        padding: const EdgeInsets.all(8.0),
                        child: TextField(
                          onChanged: (value) {
                            filterSearchResults(value);
                          },
                          controller: editingController,
                          decoration: InputDecoration(
                              labelText: "Search",
                              hintText: "Search",
                              prefixIcon: Icon(Icons.search),
                              border: OutlineInputBorder(
                                  borderRadius:
                                      BorderRadius.all(Radius.circular(10.0)))),
                        ),
                      ),
                    ],
                  ),
                ),
                SizedBox(
                  height: 400,
                  child: SingleChildScrollView(
                    child: Column(children: <Widget>[
                      ListView.builder(
                        physics: NeverScrollableScrollPhysics(),
                        shrinkWrap: true,
                        itemCount: items.length,
                        itemBuilder: (context, index) {
                          return ListTile(
                            title: Text('${items[index]}'),
                          );
                        },
                      ),
                    ]),
                  ),
                ),
              ]),
            ),
            Container(
              decoration: new BoxDecoration(color: Colors.yellow),
              child: Column(

                children: <Widget>[
                  Row(
                    children: <Widget>[
                      SpotifyButton,
                      Text(
                        "TITLE - Artist",
                        style: TextStyle(color: Colors.black),

                      ),
                    ],
                  ),
                  Slider(
                      onChanged: (duration) {},
                      value: 10,
                      max: 100,
                      min: 0,
                      activeColor: Colors.black),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: <Widget>[
                      IconButton(
                        onPressed: () {},
                        icon: Icon(FontAwesomeIcons.backward),
                        color: Colors.black,
                      ),
                      IconButton(
                        icon: Icon(
                          _isPlaying
                              ? FontAwesomeIcons.pause
                              : FontAwesomeIcons.play,
                          color: Colors.black,
                        ),
                        onPressed: () {
                          if (_isPlaying) {
                            _audioPlayer.pause();
                            setState(() {
                              _isPlaying = false;
                            });
                          } else {
                            _audioPlayer.resume();
                            setState(() {
                              _isPlaying = true;
                            });
                          }
                        },
                      ),
                      IconButton(
                        onPressed: () {},
                        icon: Icon(FontAwesomeIcons.forward),
                        color: Colors.black,
                      ),
                    ],
                  ),
                ],
              ),
            ),

          ],
        ),
      ),
    );
  }
}
