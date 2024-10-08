import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;
import 'login_page.dart';

class ProfilePage extends StatefulWidget {
  @override
  _ProfilePageState createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  String? token;
  Map<String, dynamic>? userData;
  bool isLoading = true;
  final int userId = 1; // ID de l'utilisateur à récupérer, tu peux le récupérer dynamiquement si besoin

  @override
  void initState() {
    super.initState();
    _loadToken();
  }

  Future<void> _loadToken() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    String? storedToken = prefs.getString('token');
    setState(() {
      token = storedToken;
    });
    if (token != null) {
      _fetchUserProfile();
    }
  }

  Future<void> _fetchUserProfile() async {
    final response = await http.get(
      Uri.parse('http://localhost:8080/api/users/$userId'),
      headers: {
        'Authorization': 'Bearer $token',
      },
    );

    if (response.statusCode == 200) {
      setState(() {
        userData = json.decode(response.body);
        isLoading = false;
      });
    } else {
      setState(() {
        isLoading = false;
      });
      // Gérer l'erreur de récupération des données utilisateur
      print('Erreur lors du chargement des données utilisateur');
    }
  }

  Future<void> _logout(BuildContext context) async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.remove('token'); // Supprime le token pour déconnecter l'utilisateur

    // Redirige vers la page de connexion
    Navigator.pushAndRemoveUntil(
      context,
      MaterialPageRoute(builder: (context) => LoginPage()),
      (Route<dynamic> route) => false,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Profil'),
        actions: [
          IconButton(
            icon: Icon(Icons.logout),
            onPressed: () => _logout(context), // Déconnexion au clic
          )
        ],
      ),
      body: isLoading
          ? Center(child: CircularProgressIndicator())
          : userData != null
              ? Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: <Widget>[
                      CircleAvatar(
                        radius: 50,
                        backgroundImage: NetworkImage(
                            userData!['picture'] ?? 'https://via.placeholder.com/150'),
                      ),
                      SizedBox(height: 20),
                      Text(
                        'Nom: ${userData!['lastname']}',
                        style: TextStyle(fontSize: 20),
                      ),
                      SizedBox(height: 10),
                      Text(
                        'Prénom: ${userData!['firstname']}',
                        style: TextStyle(fontSize: 20),
                      ),
                      SizedBox(height: 10),
                      Text(
                        'Email: ${userData!['email']}',
                        style: TextStyle(fontSize: 18),
                      ),
                      SizedBox(height: 20),
                      ElevatedButton(
                        onPressed: () => _logout(context),
                        child: Text('Se déconnecter'),
                      ),
                    ],
                  ),
                )
              : Center(
                  child: Text('Erreur de chargement des données utilisateur'),
                ),
    );
  }
}
