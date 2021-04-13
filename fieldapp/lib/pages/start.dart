import 'dart:async';

import 'package:fieldapp/auth.dart';
import 'package:fieldapp/pages/login.dart';
import 'package:fieldapp/types/types.dart';
import 'package:flutter/services.dart';
import 'package:fieldapp/main.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:page_transition/page_transition.dart';

class StartPage extends StatefulWidget {
  final bool transition;
  StartPage({this.transition});

  @override
  StartPageState createState() => StartPageState();
}

class StartPageState extends State<StartPage> with TickerProviderStateMixin {
  bool start = false;

  static bool attemptedAutoLogin = false;

  @override
  void initState() {
    super.initState();

    // being used as login transition animation
    if (widget.transition != null && me != null) {
      start = true;

      SystemChrome.setSystemUIOverlayStyle(SystemUiOverlayStyle(
        statusBarIconBrightness: Brightness.light,
      ));

      Timer(const Duration(milliseconds: 1500), () {
        navKey.currentState.pushReplacementNamed("/menu");
      });
      return;
    }

    // intro splash animation
    Timer(const Duration(seconds: 1), () {
      setState(() => start = true);
      SystemChrome.setSystemUIOverlayStyle(SystemUiOverlayStyle(
        statusBarIconBrightness: Brightness.light,
      ));

      Timer(const Duration(seconds: 2), () async {
        SystemChrome.setSystemUIOverlayStyle(SystemUiOverlayStyle(
          statusBarIconBrightness: Brightness.dark,
        ));

        // Attempt login
        if (await attemptLogin()) {
          navKey.currentState.pushReplacementNamed("/menu");
          return;
        }

        // Clear token (if it exists)
        Auth.signOut();

        navKey.currentState.pushReplacement(
          PageTransition(
            type: PageTransitionType.fade,
            duration: Duration(milliseconds: 1500),
            child: LoginPage(),
          ),
        );
      });
    });
  }

  Future<bool> attemptLogin() async {
    if (attemptedAutoLogin) return false;
    attemptedAutoLogin = true;

    if (me != null) return true; // already logged in

    // Get user
    GQLResponse response = await Auth.getMe();
    return response.success;
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: () async => false,
      child: Scaffold(
        body: Stack(
          children: <Widget>[
            AnimatedSize(
              duration: Duration(seconds: 1),
              vsync: this,
              curve: Curves.easeInOut,
              child: Hero(
                tag: "login-overlay",
                child: Container(
                  height: start ? double.infinity : 0,
                  padding: EdgeInsets.only(left: 40, right: 40),
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      stops: [0, 1],
                      colors: [
                        COLOUR_PRIMARY,
                        COLOUR_SECONDARY,
                      ],
                    ),
                  ),
                ),
              ),
            ),
            Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: <Widget>[
                  Hero(
                    tag: "logo",
                    child: AnimatedSwitcher(
                      duration: Duration(seconds: 1),
                      child: Container(
                        key: ValueKey<String>("logo-$start"),
                        child: FaIcon(
                          FontAwesomeIcons.solidSteak,
                          color: start ? Colors.white : COLOUR_PRIMARY,
                          size: 80,
                        ),
                      ),
                    ),
                  ),
                  SizedBox(height: 34),
                  Hero(
                    tag: "title",
                    child: Material(
                      type: MaterialType.transparency,
                      child: AnimatedSwitcher(
                        duration: Duration(seconds: 1),
                        child: Container(
                          key: ValueKey<String>("title-$start"),
                          width: 200,
                          child: Text(
                            "Genesis",
                            style: TextStyle(
                              color: start ? Colors.white : COLOUR_PRIMARY,
                              fontWeight: FontWeight.w600,
                              fontSize: 48,
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

Widget startOverlayAnim(
    {AlignmentGeometry alignment = Alignment.bottomCenter, double height = 0}) {
  return Align(
    alignment: alignment,
    child: Hero(
      tag: "login-overlay",
      child: Container(
        height: height,
        padding: EdgeInsets.only(left: 40, right: 40),
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            stops: [0, 1],
            colors: [
              COLOUR_PRIMARY,
              COLOUR_SECONDARY,
            ],
          ),
        ),
      ),
    ),
  );
}
