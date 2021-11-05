import 'package:flutter/material.dart';
import 'package:musicroomapp/style.dart';

class MyPlaylist extends StatelessWidget {
  const MyPlaylist({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Container(
          padding: const EdgeInsets.all(80.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                'Welcome to Music Room',
                style: Theme.of(context).textTheme.headline1,
              ),
              Text(
                'Current Room available',
                style: Theme.of(context).textTheme.headline1,
              ),
              ElevatedButton(
                child: const Text('EXIT'),
                onPressed: () {
                  Navigator.pushReplacementNamed(context, '/');
                },
                style: ElevatedButton.styleFrom(
                  primary: Colors.yellow,
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}