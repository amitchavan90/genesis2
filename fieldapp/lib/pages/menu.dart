import 'package:fieldapp/main.dart';
import 'package:fieldapp/types/enums.dart';
import 'package:flutter/material.dart';
import 'package:pinput/pin_put/pin_put.dart';

const String PIN = "1234";

class HomeOptionsPage extends StatefulWidget {
  HomeOptionsPage({Key key}) : super(key: key);

  @override
  _HomeOptionsPageState createState() => _HomeOptionsPageState();
}

class _HomeOptionsPageState extends State<HomeOptionsPage> {
  final TextEditingController _pinPutController = TextEditingController();
  FocusNode _pinPutFocusNode = FocusNode();

  // pin code input
  Widget pinPut() {
    BoxDecoration pinPutDecoration = BoxDecoration(
      color: Colors.white,
      borderRadius: BorderRadius.circular(0),
      border: Border.all(
        color: COLOUR_PRIMARY,
        width: 2.0,
      ),
    );
    return Container(
        padding: EdgeInsets.symmetric(horizontal: 15),
        child: Column(
          children: <Widget>[
            Text(
              "Enter Pin:",
              style: TextStyle(fontSize: 25),
            ),
            Container(
              margin: EdgeInsets.symmetric(vertical: 10),
              child: Text(
                "(tap one of the boxes to enter)",
                style: TextStyle(
                  fontSize: 15,
                ),
              ),
            ),
            PinPut(
              preFilledChar: "",
              eachFieldWidth: 50,
              eachFieldHeight: 60,
              fieldsCount: 4,
              focusNode: _pinPutFocusNode,
              controller: _pinPutController,
              onSubmit: (String pin) {},
              submittedFieldDecoration: pinPutDecoration,
              selectedFieldDecoration: pinPutDecoration,
              followingFieldDecoration: pinPutDecoration,
              pinAnimationType: PinAnimationType.scale,
              textStyle: TextStyle(color: COLOUR_PRIMARY, fontSize: 20),
            ),
          ],
        ));
  }

  Widget build(BuildContext context) {
    bool isAdvancedUser = me.role.tier == 1 ||
        me.role.permissions.indexOf(Perm.UseAdvancedMode) != -1;
    return Scaffold(
      body: SingleChildScrollView(
        child: SafeArea(
          child: Center(
            child: Container(
              margin: const EdgeInsets.symmetric(horizontal: 40),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: <Widget>[
                  Container(
                    child: Text(
                      "",
                      style: TextStyle(color: COLOUR_PRIMARY, fontSize: 40),
                    ),
                  ),
                  Container(
                    width: double.infinity,
                    child: RaisedButton(
                      padding: EdgeInsets.symmetric(vertical: 35),
                      color: COLOUR_PRIMARY,
                      child: Text(
                        "Pack Carton and Complete Task",
                        style: TextStyle(color: Colors.white),
                      ),
                      onPressed: () {
                        mode = Mode.packCartonCompleteTask;
                        navKey.currentState.pushReplacementNamed("/home");
                      },
                    ),
                  ),
                  Container(
                    margin: const EdgeInsets.symmetric(vertical: 15),
                    width: double.infinity,
                    child: RaisedButton(
                      padding: EdgeInsets.symmetric(vertical: 35),
                      color: COLOUR_PRIMARY,
                      child: Text(
                        "Complete Task",
                        style: TextStyle(color: Colors.white),
                      ),
                      onPressed: () {
                        mode = Mode.completeTask;
                        navKey.currentState.pushReplacementNamed("/home");
                      },
                    ),
                  ),
                  Container(
                    margin: const EdgeInsets.only(top: 90, bottom: 15),
                    child: RaisedButton(
                      padding: EdgeInsets.all(25),
                      color: COLOUR_PRIMARY,
                      disabledColor: Colors.grey,
                      child: Text(
                        "Advanced User",
                        style: TextStyle(color: Colors.white),
                      ),
                      onPressed: isAdvancedUser
                          ? () {
                              mode = Mode.advance;
                              navKey.currentState.pushReplacementNamed("/home");
                            }
                          : null,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
