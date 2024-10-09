import 'package:flutter/material.dart';
import 'models/kermesse.dart'; // Le modèle que tu as créé
import 'services/kermesse_service.dart'; // Le service pour récupérer les kermesses

class KermesseListPage extends StatefulWidget {
  @override
  _KermesseListPageState createState() => _KermesseListPageState();
}

class _KermesseListPageState extends State<KermesseListPage> {
  late Future<List<Kermesse>> _futureKermesses;

  @override
  void initState() {
    super.initState();
    _futureKermesses = KermesseService().fetchKermesses();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Liste des Kermesses'),
      ),
      body: FutureBuilder<List<Kermesse>>(
        future: _futureKermesses,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Erreur : ${snapshot.error}'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return Center(child: Text('Aucune kermesse trouvée.'));
          } else {
            return ListView.builder(
              itemCount: snapshot.data!.length,
              itemBuilder: (context, index) {
                final kermesse = snapshot.data![index];
                return Card(
                  child: ListTile(
                    title: Text(kermesse.name),
                    subtitle: Text(kermesse.description),
                    trailing: Text(kermesse.date.toLocal().toString()),
                    onTap: () {
                      // Tu peux ajouter une action pour chaque kermesse, par exemple pour afficher plus de détails
                    },
                  ),
                );
              },
            );
          }
        },
      ),
    );
  }
}
