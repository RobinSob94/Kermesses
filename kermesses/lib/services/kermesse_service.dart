import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/kermesse.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

class KermesseService {
  final String apiUrl = '${dotenv.env['API_PROTOCOL']}://${dotenv.env['API_HOST']}:${dotenv.env['API_PORT']}/kermesses'; // Ton endpoint pour récupérer les kermesses

  Future<List<Kermesse>> fetchKermesses() async {
    final response = await http.get(Uri.parse(apiUrl));

    if (response.statusCode == 200) {
      List<dynamic> body = jsonDecode(response.body);
      List<Kermesse> kermesses = body.map((dynamic item) => Kermesse.fromJson(item)).toList();
      return kermesses;
    } else {
      throw Exception('Failed to load kermesses');
    }
  }
}
