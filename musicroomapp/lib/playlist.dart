import 'package:flutter/material.dart';
import 'package:musicroomapp/style.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:audioplayers/audioplayers.dart';
import 'package:file_picker/file_picker.dart';

var SpotifyIcon = Icon(
  FontAwesomeIcons.spotify,
  color: Colors.black,
);
// in case we need to sync something with spotify first
var SpotifyButton = IconButton(icon: SpotifyIcon, onPressed: () {});

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
