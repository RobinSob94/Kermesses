import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:kermesses/main.dart';
import 'package:kermesses/login_page.dart';
import 'package:kermesses/profile_page.dart';

void main() {
  testWidgets('Login page loads and shows necessary fields', (WidgetTester tester) async {
    // Build the login page and trigger a frame.
    await tester.pumpWidget(MyApp());

    // Verify that the login page contains email and password fields and a login button.
    expect(find.byType(TextFormField), findsNWidgets(2)); // Deux champs de texte (email, mot de passe)
    expect(find.byType(ElevatedButton), findsOneWidget); // Bouton de connexion

    // Verify that the login button has the correct text.
    expect(find.text('Se connecter'), findsOneWidget);
  });

  testWidgets('Profile page displays user information', (WidgetTester tester) async {
    // Mock some user data for testing purposes
    final Map<String, dynamic> mockUserData = {
      'firstname': 'John',
      'lastname': 'Doe',
      'email': 'john.doe@example.com',
      'picture': 'https://via.placeholder.com/150',
    };

    // Simulate that the user is logged in and the profile page loads.
    await tester.pumpWidget(MaterialApp(
      home: ProfilePage(),
    ));

    // Pump again to ensure the page is fully loaded.
    await tester.pump();

    // Check if the user's information is displayed correctly.
    expect(find.text('Nom: Doe'), findsOneWidget);
    expect(find.text('Prénom: John'), findsOneWidget);
    expect(find.text('Email: john.doe@example.com'), findsOneWidget);
    expect(find.byType(CircleAvatar), findsOneWidget); // Image de l'utilisateur

    // Check if the logout button is present.
    expect(find.byType(ElevatedButton), findsOneWidget);
    expect(find.text('Se déconnecter'), findsOneWidget);
  });

  testWidgets('Login and logout flow works', (WidgetTester tester) async {
    // Start with the login page.
    await tester.pumpWidget(MyApp());

    // Enter email and password
    await tester.enterText(find.byType(TextFormField).at(0), 'john.doe@example.com');
    await tester.enterText(find.byType(TextFormField).at(1), 'password123');

    // Tap the login button.
    await tester.tap(find.text('Se connecter'));
    await tester.pump();

    // Mock the login success and navigate to the profile page.
    await tester.pumpWidget(MaterialApp(
      home: ProfilePage(),
    ));
    await tester.pump();

    // Ensure the profile page is displayed.
    expect(find.text('Nom: Doe'), findsOneWidget);
    expect(find.text('Prénom: John'), findsOneWidget);

    // Tap the logout button.
    await tester.tap(find.text('Se déconnecter'));
    await tester.pump();

    // Verify that after logout, the login page is shown again.
    expect(find.byType(LoginPage), findsOneWidget);
  });
}
