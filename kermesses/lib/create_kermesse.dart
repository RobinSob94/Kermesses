import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:jwt_decode/jwt_decode.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';



class CreateKermessePage extends StatefulWidget {
  @override
  _CreateKermessePageState createState() => _CreateKermessePageState();
}

class _CreateKermessePageState extends State<CreateKermessePage> {
  final _formKey = GlobalKey<FormState>();
  String _kermesseName = '';
  bool _isLoading = false;

  // Fonction pour récupérer le token depuis SharedPreferences
  Future<String?> _getToken() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    return prefs.getString('token');
  }

  // Fonction pour vérifier si l'utilisateur a les rôles requis
  bool _userHasPermission(String token) {
    Map<String, dynamic> payload = Jwt.parseJwt(token);
    int userRole = payload['role'];
    return userRole == 1 || userRole == 2;
  }

  // Fonction pour envoyer les données du formulaire à l'API
  Future<void> _createKermesse() async {
    if (_formKey.currentState!.validate()) {
      setState(() {
        _isLoading = true;
      });
      _formKey.currentState!.save();

      String? token = await _getToken();

      if (token != null) {
        if (!_userHasPermission(token)) {
          ScaffoldMessenger.of(context).showSnackBar(SnackBar(
            content: Text("Vous n'avez pas l'autorisation de créer une kermesse."),
          ));
          setState(() {
            _isLoading = false;
          });
          return;
        }

        try {
          var response = await http.post(
            Uri.parse('${dotenv.env['API_PROTOCOL']}://${dotenv.env['API_HOST']}:${dotenv.env['API_PORT']}/create-kermesse'),
            headers: {
              'Content-Type': 'application/json',
              'Authorization': 'Bearer $token',
            },
            body: json.encode({
              'name': _kermesseName,
            }),
          );

          if (response.statusCode == 201) {
            ScaffoldMessenger.of(context).showSnackBar(SnackBar(
              content: Text("Kermesse créée avec succès !"),
            ));
          } else {
            ScaffoldMessenger.of(context).showSnackBar(SnackBar(
              content: Text("Erreur lors de la création de la kermesse : ${response.body}"),
            ));
          }
        } catch (e) {
          ScaffoldMessenger.of(context).showSnackBar(SnackBar(
            content: Text("Erreur réseau : $e"),
          ));
        }

        setState(() {
          _isLoading = false;
        });
      } else {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text("Utilisateur non connecté."),
        ));
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Créer une Kermesse'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                decoration: InputDecoration(labelText: 'Nom de la Kermesse'),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Veuillez entrer un nom de kermesse';
                  }
                  return null;
                },
                onSaved: (value) {
                  _kermesseName = value!;
                },
              ),
              SizedBox(height: 20),
              _isLoading
                  ? CircularProgressIndicator()
                  : ElevatedButton(
                      onPressed: _createKermesse,
                      child: Text('Créer Kermesse'),
                    ),
            ],
          ),
        ),
      ),
    );
  }
}
