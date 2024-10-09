import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'login_page.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'profile_page.dart';

void main() async {
await dotenv.load(fileName: '.env');
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Kermesses App',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: CheckAuth(), // Vérification de l'état d'authentification
    );
  }
}

class CheckAuth extends StatefulWidget {
  @override
  _CheckAuthState createState() => _CheckAuthState();
}

class _CheckAuthState extends State<CheckAuth> {
  String? token;

  @override
  void initState() {
    super.initState();
    _checkLoginStatus();
  }

  // Vérification du statut de connexion
  _checkLoginStatus() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    String? storedToken = prefs.getString('token');

    setState(() {
      token = storedToken;
    });
  }

  // Déconnexion
  _logout() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.remove('token');
    setState(() {
      token = null; // Reset du token
    });
  }

  @override
  Widget build(BuildContext context) {
    if (token != null) {
      return ProfilePage(logoutCallback: _logout); // Redirige vers la page de profil
    } else {
      return LoginPage(); // Redirige vers la page de connexion
    }
  }
}
