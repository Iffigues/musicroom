import 'package:flutter/material.dart';
import 'package:musicroomapp/login.dart';
import 'package:musicroomapp/signin.dart';
import 'package:musicroomapp/style.dart';
import 'package:musicroomapp/playlist.dart';
import 'package:musicroomapp/listmodel.dart';

import 'package:provider/provider.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        Provider(create: (context) => ListModel()),
      ],
      child: MaterialApp(
        title: 'Provider Demo',
        theme: appTheme,
        initialRoute: '/',
        routes: {
          '/': (context) => const MyLogin(),
          '/signin': (context) => const MySignIn(),
          '/playlist': (context) => const MyPlaylist(),
        },
      ),
    );
  }
}
