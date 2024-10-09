import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'dart:convert'; // Pour décoder la réponse JSON
import 'package:http/http.dart' as http;
import 'package:jwt_decode/jwt_decode.dart'; // Pour décoder le token JWT
import 'package:flutter_dotenv/flutter_dotenv.dart';


class ProfilePage extends StatefulWidget {
  final VoidCallback logoutCallback;

  ProfilePage({required this.logoutCallback});

  @override
  _ProfilePageState createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  Map<String, dynamic>? userData;
  bool isLoading = true;
  String? userId;

  @override
  void initState() {
    super.initState();
    _fetchUserProfile();
  }

  // Récupérer les informations de l'utilisateur à partir de l'API
  _fetchUserProfile() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    String? token = prefs.getString('token');

    if (token != null) {
      // Décoder le token pour obtenir l'ID utilisateur
      Map<String, dynamic> decodedToken = Jwt.parseJwt(token);
      userId = decodedToken['id'].toString(); // Récupère l'ID de l'utilisateur

      // Appel API pour récupérer les informations de l'utilisateur à partir de l'ID
      final response = await http.get(
        Uri.parse('${dotenv.env['API_PROTOCOL']}://${dotenv.env['API_HOST']}:${dotenv.env['API_PORT']}/api/users/$userId'),
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
        // Gestion d'erreur si la requête échoue
        setState(() {
          isLoading = false;
        });
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur lors du chargement du profil.')),
        );
      }
    }
  }

  // Déconnexion
  _logout() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.remove('token');
    widget.logoutCallback(); // Déclenche la déconnexion et redirige
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Mon Profil"),
        actions: [
          IconButton(
            icon: Icon(Icons.logout),
            onPressed: _logout, // Bouton de déconnexion
          ),
        ],
      ),
      body: isLoading
          ? Center(child: CircularProgressIndicator()) // Affiche un loader si les données sont en cours de chargement
          : userData != null
              ? Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        "Nom : ${userData!['name']}",
                        style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                      ),
                      SizedBox(height: 8),
                      Text(
                        "Email : ${userData!['email']}",
                        style: TextStyle(fontSize: 18),
                      ),
                      SizedBox(height: 16),
                      Text(
                        "Rôle : ${userData!['role']}",
                        style: TextStyle(fontSize: 18),
                      ),
                      SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _logout,
                        child: Text('Déconnexion'),
                      ),
                    ],
                  ),
                )
              : Center(child: Text("Erreur lors du chargement du profil.")),
    );
  }
}
