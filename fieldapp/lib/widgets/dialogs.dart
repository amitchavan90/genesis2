import 'package:awesome_dialog/anims/anims.dart';
import 'package:awesome_dialog/awesome_dialog.dart';
import 'package:fieldapp/main.dart';
import 'package:fieldapp/types/types.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/widgets.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

void openListDialog(
    {@required BuildContext context,
    @required String title,
    IconData icon,
    List<Widget> children,
    bool slideFromRight = false}) {
  List<Widget> c = [
    Container(
      decoration: BoxDecoration(color: COLOUR_PRIMARY),
      padding: EdgeInsets.symmetric(vertical: 12, horizontal: 20),
      child: Row(
        children: <Widget>[
          icon != null
              ? Padding(
                  child: FaIcon(
                    icon,
                    color: Colors.white,
                  ),
                  padding: EdgeInsets.only(right: 10),
                )
              : Container(),
          Text(
            title,
            style: TextStyle(
              color: Colors.white,
              fontSize: 20,
            ),
          ),
        ],
      ),
    ),
  ];
  if (children != null) c.addAll(children);

  showDialog<void>(
    context: context,
    barrierDismissible: true,
    builder: (BuildContext context) {
      return Slide(
        from: slideFromRight ? SlideFrom.RIGHT : SlideFrom.LEFT,
        slideDistance: MediaQuery.of(context).size.width,
        duration: 0.4,
        child: Center(
          child: Material(
            type: MaterialType.transparency,
            child: Container(
              decoration: BoxDecoration(color: Colors.white),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: c,
              ),
            ),
          ),
        ),
      );
    },
  );
}

Future<void> showErrorDialog(
  BuildContext context,
  String message,
) {
  return AwesomeDialog(
    context: context,
    dialogType: DialogType.ERROR,
    animType: AnimType.BOTTOMSLIDE,
    tittle: 'Error',
    desc: message,
    btnCancelOnPress: () {},
    btnCancelColor: COLOUR_PRIMARY,
    btnCancelText: "Ok",
    headerAnimationLoop: false,
  ).show();
}

Future<void> showSuccessDialog(
  BuildContext context,
  String title,
  String message,
) {
  return AwesomeDialog(
    context: context,
    dialogType: DialogType.SUCCES,
    animType: AnimType.BOTTOMSLIDE,
    tittle: title,
    desc: message,
    btnCancelOnPress: () {},
    btnCancelColor: COLOUR_PRIMARY,
    btnCancelText: "Ok",
    headerAnimationLoop: false,
  ).show();
}

Future<void> showWarningDialog(
  BuildContext context,
  String title,
  String message,
) {
  return AwesomeDialog(
    context: context,
    dialogType: DialogType.WARNING,
    animType: AnimType.BOTTOMSLIDE,
    tittle: title,
    desc: message,
    btnCancelOnPress: () {},
    btnCancelColor: COLOUR_PRIMARY,
    btnCancelText: "Ok",
    headerAnimationLoop: false,
  ).show();
}

Future<void> showInfoDialog(
    BuildContext context, String title, String message, Function onPressed) {
  return AwesomeDialog(
    context: context,
    dialogType: DialogType.INFO,
    animType: AnimType.BOTTOMSLIDE,
    tittle: title,
    desc: message,
    btnCancelColor: COLOUR_PRIMARY,
    btnCancelText: "Ok",
    btnOkOnPress: () {
      onPressed();
    },
    headerAnimationLoop: false,
  ).show();
}

Future<void> showScanNextDialog(
    BuildContext context, String title, String message, Function onPressed) {
  return AwesomeDialog(
    context: context,
    dialogType: DialogType.INFO,
    animType: AnimType.BOTTOMSLIDE,
    tittle: title,
    desc: message,
    btnCancelColor: COLOUR_PRIMARY,
    btnCancelText: "Stop Scan",
    btnOkText: "Scan Next",
    btnOkOnPress: () {
      onPressed();
    },
    btnCancelOnPress: () {
      return;
    },
    headerAnimationLoop: false,
  ).show();
}

Future<void> showEnterCodeDialog(
  BuildContext context,
  GlobalKey<FormState> formKey,
  Future<bool> Function(String) verify, {
  String msg,
  Function onCancel,
  bool barrierDismissible = false,
  Future<void> Function() resend,
}) {
  TextEditingController inputController = TextEditingController();

  List<Widget> actions = [
    FlatButton(
      child: Text("Verify"),
      onPressed: () async {
        if (!formKey.currentState.validate()) return;
        if (await verify(inputController.text)) {
          Navigator.of(context).pop();
        }
      },
      color: Colors.grey.withOpacity(0.3),
      padding: EdgeInsets.only(left: 30, right: 30),
    ),
    FlatButton(
      child: Text("Cancel"),
      onPressed: () {
        if (onCancel != null) onCancel();
        Navigator.of(context).pop();
      },
    ),
  ];

  if (resend != null) {
    actions.insert(
      1,
      FlatButton(
        child: Text("Resend"),
        onPressed: resend,
      ),
    );
  }

  return showDialog<void>(
    context: context,
    barrierDismissible: barrierDismissible,
    builder: (BuildContext context) {
      return AlertDialog(
        title: Text("Enter Code"),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: <Widget>[
            msg != null ? Text(msg) : Container(),
            Form(
              key: formKey,
              child: TextFormField(
                keyboardType: TextInputType.text,
                autofocus: true,
                maxLength: 6,
                maxLengthEnforced: true,
                controller: inputController,
                textCapitalization: TextCapitalization.characters,
                validator: (value) {
                  if (value.length < 6) return ' A valid code is required ';
                  return null;
                },
                onFieldSubmitted: (value) async {
                  if (await verify(value)) Navigator.of(context).pop();
                },
                decoration: InputDecoration(
                  labelText: "Code",
                  hintText: '...',
                ),
              ),
            ),
          ],
        ),
        actions: actions,
        scrollable: true,
        contentPadding: EdgeInsets.symmetric(vertical: 5, horizontal: 15),
      );
    },
  );
}

Future<void> captureCartonDialog(
  dynamic carton,
  BuildContext context,
  Future<void> Function(dynamic) captureCartonImage,
  TrackAction selectedTrackAction,
  Function(bool) setLoading,
) async {
  if (selectedTrackAction == null ||
      selectedTrackAction.requirePhotos.length == 0 ||
      !selectedTrackAction.requirePhotos[0]) {
    return;
  }

  // if already has photo
  if (mode != Mode.advance) {
    if (carton["cartonPhoto"] != null) {
      showInfoDialog(context, "This carton already has a photo", "", () {});
      return;
    }
    // doesnt have photo
    AwesomeDialog(
      dismissOnTouchOutside: false,
      headerAnimationLoop: false,
      context: context,
      dialogType: DialogType.INFO,
      animType: AnimType.BOTTOMSLIDE,
      body: WillPopScope(
        onWillPop: () async => false,
        child: Container(
          transform: Matrix4.translationValues(0.0, -40.0, 0.0),
          child: SingleChildScrollView(
            child: Container(
              child: Center(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: <Widget>[
                    Container(
                      child: Align(
                        child: CloseButton(
                            color: Colors.black,
                            onPressed: () {
                              Navigator.pop(context);
                            }),
                        alignment: Alignment.topRight,
                      ),
                    ),
                    Center(
                      child: Text(
                        "Take Photo of Carton Label next to QR. Example:",
                        style: TextStyle(fontSize: 18),
                      ),
                    ),
                    SizedBox(
                      height: 10,
                    ),
                    Image.asset('assets/cartonAndQR.png')
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
      desc: "",
      btnOkOnPress: () async {
        // brings up camera
        await captureCartonImage(carton);
      },
    ).show();
    return;
  }
  await captureCartonImage(carton);
}

Future<void> captureProductDialog(
    dynamic carton,
    BuildContext context,
    Future<void> Function(dynamic) captureProductImage,
    TrackAction selectedTrackAction,
    Function(bool) setLoading) async {
  if (selectedTrackAction == null ||
      selectedTrackAction.requirePhotos.length < 2 ||
      !selectedTrackAction.requirePhotos[1]) {
    return;
  }

  // if already has photo
  if (mode != Mode.advance) {
    if (carton["productPhoto"] != null) {
      showInfoDialog(
          context, "This carton already has a photo of a product", "", () {});
      return;
    }

    // doesnt have photo
    AwesomeDialog(
      dismissOnTouchOutside: false,
      headerAnimationLoop: false,
      context: context,
      dialogType: DialogType.INFO,
      animType: AnimType.BOTTOMSLIDE,
      body: WillPopScope(
        onWillPop: () async => false,
        child: Container(
          transform: Matrix4.translationValues(0.0, -40.0, 0.0),
          child: Container(
            child: Center(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: <Widget>[
                  Container(
                    child: Align(
                      child: CloseButton(
                          color: Colors.black,
                          onPressed: () {
                            Navigator.pop(context);
                          }),
                      alignment: Alignment.topRight,
                    ),
                  ),
                  Center(
                    child: Text(
                      "Take Photo of a Product from the Carton next to QR. Example:",
                      style: TextStyle(fontSize: 18),
                    ),
                  ),
                  SizedBox(
                    height: 10,
                  ),
                  Image.asset('assets/productAndQR.png')
                ],
              ),
            ),
          ),
        ),
      ),
      desc: "",
      btnOkOnPress: () {
        // brings up camera
        captureProductImage(carton);
      },
    ).show();
    return;
  }
  await captureProductImage(carton);
}
