import 'package:flutter/material.dart';
import 'package:form_field_validator/form_field_validator.dart';
import 'package:http/http.dart' as http;
class MySignIn extends StatefulWidget {
  const MySignIn({Key? key}) : super(key: key);

  @override
  MySignInFormState createState() {
    return MySignInFormState();
  }
}

class MySignInFormState extends State<MySignIn> {
  final _formKey = GlobalKey<FormState>();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Container(
          padding: const EdgeInsets.fromLTRB(80.0, 20.0, 80.0, 20.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              ElevatedButton(
                child: const Text('LOGIN'),
                onPressed: () {
                  Navigator.pushReplacementNamed(context, '/');
                },
                style: ElevatedButton.styleFrom(
                  primary: Colors.yellow,
                ),
              ),
              Text(
                'Sign in',
                style: Theme.of(context).textTheme.headline1,
              ),
              Form(
                key: _formKey,
                child: Column(
                  children: [
                    TextFormField(
                        decoration: const InputDecoration(
                          hintText: 'Email',
                        ),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Enter an email address';
                          }
                          String pattern = r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$';
                          RegExp regExp = RegExp(pattern);
                          if (!regExp.hasMatch(value)) {
                            return 'Enter a valid email address';
                          }
                        }),
                    TextFormField(
                        decoration: const InputDecoration(
                          hintText: 'Username',
                        ),
                        validator: RequiredValidator(
                            errorText: 'A username is required')),
                    TextFormField(
                        decoration: const InputDecoration(
                          hintText: 'Password',
                        ),
                        obscureText: true,
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'enter a password';
                          }
                          if (value.length < 8) {
                            return 'password must contain at least 8 characters';
                          }
                          String pattern =
                              r'^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[!@#\$&*~]).{8,}$';
                          RegExp regExp = RegExp(pattern);
                          if (!regExp.hasMatch(value)) {
                            return 'password must contain at least 1 special character, 1 lowcase, 1 upcase, 1 digit';
                          }
                        }),
                    const SizedBox(
                      height: 24,
                    ),
                  ],
                ),
              ),
              ElevatedButton(
                child: const Text('REGISTER'),
                onPressed: () {
                  if (_formKey.currentState!.validate()) {
                    Navigator.pushReplacementNamed(context, '/playlist');
                  }
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
