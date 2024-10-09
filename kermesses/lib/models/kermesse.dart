class Kermesse {
  final int id;
  final String name;
  final String description;
  final DateTime date;

  Kermesse({
    required this.id,
    required this.name,
    required this.description,
    required this.date,
  });

  factory Kermesse.fromJson(Map<String, dynamic> json) {
    return Kermesse(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      date: DateTime.parse(json['date']),
    );
  }
}
