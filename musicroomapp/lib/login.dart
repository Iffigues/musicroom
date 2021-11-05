import 'package:flutter/material.dart';
import 'package:musicroomapp/style.dart';
import 'package:form_field_validator/form_field_validator.dart';

class MyLogin extends StatefulWidget {
  const MyLogin({Key? key}) : super(key: key);

  @override
  MyCustomFormState createState() {
    return MyCustomFormState();
  }
}

class MyCustomFormState extends State<MyLogin> {
  final _formKey = GlobalKey<FormState>();

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
              )
            ],
          ),
        ),
      ),
    );
  }
}
