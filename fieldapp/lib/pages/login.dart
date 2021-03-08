import 'package:fieldapp/auth.dart';
import 'package:fieldapp/types/types.dart';
import 'package:fieldapp/widgets/changePassword.dart';
import 'package:flutter/painting.dart';
import 'package:flutter/widgets.dart';
import 'package:fieldapp/widgets/commonWidgets.dart';
import 'package:fieldapp/widgets/dialogs.dart';
import 'package:fieldapp/main.dart';
import 'package:fieldapp/pages/start.dart';
import 'package:fieldapp/utils.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:package_info/package_info.dart';
import 'package:page_transition/page_transition.dart';
import 'package:fieldapp/widgets/versionCheck.dart';
import 'package:geolocator/geolocator.dart';
import 'package:geocoder/geocoder.dart';

class LoginPage extends StatefulWidget {
  final String email;
  final String password;
  LoginPage({this.email, this.password});

  @override
  LoginPageState createState() => LoginPageState();
}

class LoginPageState extends State<LoginPage> {
  final loginFormKey = GlobalKey<FormState>();
  final verifyFormKey = GlobalKey<FormState>();
  final emailInput = TextEditingController();
  final passwordInput = TextEditingController();
  final passwordInputFocus = FocusNode();
  bool autovalidate = false;
  bool loading = false;

  String hostDropDownItem = hostOption;
  final apiHostInput = TextEditingController(text: host);

  String version = "";

  @override
  void initState() {
    super.initState();

    if (widget.email != null) emailInput.text = widget.email;
    if (widget.password != null) passwordInput.text = widget.password;

    getVersion();
    versionCheck(context);

    if (!prefs.containsKey("hostOption")) setDefaultAPIHost();
  }

  void getVersion() async {
    PackageInfo packageInfo = await PackageInfo.fromPlatform();
    setState(() {
      // packageInfo.version is always 1.0.0 when debugging
      version =
          packageInfo.version == "1.0.0" ? "DEBUGGING" : packageInfo.version;
    });
  }

  @override
  void dispose() {
    emailInput?.dispose();
    passwordInput?.dispose();
    passwordInputFocus?.dispose();
    super.dispose();
  }

  void signIn() async {
    if (!loginFormKey.currentState.validate()) {
      if (autovalidate) return;
      setState(() => autovalidate = true);
      return;
    }

    FocusScope.of(context).unfocus(); // close keyboard

    // Login
    setState(() => loading = true);
    GQLResponse response =
        await Auth.signIn(emailInput.text, passwordInput.text);

    // Error check
    if (!response.success) {
      setState(() => loading = false);
      showErrorDialog(context, response.message);
      return;
    }

    // Get user
    GQLResponse response2 = await Auth.getMe();

    // Error check
    if (!response2.success) {
      setState(() => loading = false);
      showErrorDialog(context, response2.message);
      return;
    }

    setState(() => loading = false);

    navKey.currentState.pushReplacement(
      PageTransition(
        type: PageTransitionType.fade,
        duration: Duration(milliseconds: 1500),
        child: StartPage(transition: true),
      ),
    );
  }

  void onForgotPassword() async {
    setState(() => loading = true);
    GQLResponse response = await Auth.forgotPassword(emailInput.text);
    setState(() => loading = false);

    if (!response.success) {
      showErrorDialog(context, response.message);
      return;
    }

    Navigator.of(context).pop();

    // Enter code dialog
    showEnterCodeDialog(
      context,
      verifyFormKey,
      (String token) async {
        GQLResponse verifyResponse =
            await Auth.verifyResetToken(token, emailInput.text);

        if (!verifyResponse.success) {
          showErrorDialog(context, verifyResponse.message);
          return false;
        }

        // Change password dialog
        Navigator.of(context).pop();
        showDialog<void>(
          context: context,
          barrierDismissible: false,
          builder: (BuildContext context) =>
              ChangePassword(token: token, email: emailInput.text),
        );

        return false; // don't pop context (already popped)
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: () async => false,
      child: GestureDetector(
        onTap: () => FocusScope.of(context)
            .requestFocus(FocusNode()), // close keyboard on tap outside
        child: Stack(
          children: <Widget>[
            // Hero animation for start page background
            startOverlayAnim(),

            Scaffold(
              body: SingleChildScrollView(
                child: Form(
                  key: loginFormKey,
                  autovalidate: autovalidate,
                  child: Container(
                    color: Colors.white,
                    height: MediaQuery.of(context).size.height,
                    padding: EdgeInsets.only(left: 26, right: 26),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: <Widget>[
                        SizedBox(height: 40),

                        // Logo
                        Center(
                          child: Hero(
                            tag: "logo",
                            child: FaIcon(
                              FontAwesomeIcons.solidSteak,
                              color: COLOUR_PRIMARY,
                              size: 80,
                            ),
                          ),
                        ),

                        SizedBox(height: 18),

                        // Title
                        Column(
                          children: [
                            Center(
                              child: Hero(
                                tag: "title",
                                child: Material(
                                  type: MaterialType.transparency,
                                  child: Text(
                                    "Genesis",
                                    style: TextStyle(
                                      fontSize: 48,
                                      color: COLOUR_PRIMARY,
                                      fontWeight: FontWeight.w600,
                                      height: 0.8,
                                    ),
                                    maxLines: 1,
                                  ),
                                ),
                              ),
                            ),
                            Material(
                              type: MaterialType.transparency,
                              child: Text(
                                version,
                                style: TextStyle(
                                  fontSize: 14,
                                  color: Colors.grey.shade700,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                            ),
                          ],
                        ),

                        Spacer(flex: 2),

                        // Email
                        Hero(
                          tag: "email",
                          child: Material(
                            type: MaterialType.transparency,
                            child: textFieldWithIcon(
                              "Email",
                              icon: FaIcon(
                                FontAwesomeIcons.solidUser,
                                color: COLOUR_PRIMARY,
                              ),
                              keyboardType: TextInputType.emailAddress,
                              controller: emailInput,
                              textInputAction: TextInputAction.next,
                              validator: (value) {
                                if (value.isEmpty)
                                  return ' This field is required ';
                                else if (!isEmail(value))
                                  return ' Please enter a valid email address ';
                                return null;
                              },
                              focusNext: passwordInputFocus,
                              context: context,
                            ),
                          ),
                        ),
                        Spacer(),

                        // Password
                        Hero(
                          tag: "password",
                          child: Material(
                            type: MaterialType.transparency,
                            child: textFieldWithIcon(
                              "Password",
                              icon: FaIcon(
                                FontAwesomeIcons.solidLockAlt,
                                color: COLOUR_PRIMARY,
                              ),
                              keyboardType: TextInputType.text,
                              obscureText: true,
                              controller: passwordInput,
                              focusNode: passwordInputFocus,
                              validator: (value) {
                                if (value.isEmpty)
                                  return ' This field is required ';
                                return null;
                              },
                              onFieldSubmitted: (value) {
                                signIn();
                              },
                            ),
                          ),
                        ),
                        Spacer(),

                        // Login Button
                        Hero(
                          tag: "login-submit",
                          child: Material(
                            type: MaterialType.transparency,
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: <Widget>[
                                Text(
                                  "Login",
                                  style: TextStyle(
                                    fontSize: 32,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                                Container(
                                  decoration: BoxDecoration(
                                    color: COLOUR_PRIMARY,
                                    borderRadius: BorderRadius.circular(30),
                                  ),
                                  width: 55,
                                  height: 55,
                                  child: Material(
                                    type: MaterialType.transparency,
                                    child: InkWell(
                                      child: Container(
                                        width: 45,
                                        height: 45,
                                        child: Center(
                                          child: FaIcon(
                                            FontAwesomeIcons.arrowRight,
                                            color: Colors.white,
                                            size: 26,
                                          ),
                                        ),
                                      ),
                                      borderRadius: BorderRadius.circular(30),
                                      onTap: () {
                                        signIn();
                                      },
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),

                        Spacer(flex: 2),

                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: <Widget>[
                            forgotPasswordButton(),

                            // Change API Host
                            Material(
                              type: MaterialType.transparency,
                              child: DropdownButton<String>(
                                value: hostDropDownItem,
                                icon: Padding(
                                  padding: EdgeInsets.only(left: 8),
                                  child: FaIcon(
                                    FontAwesomeIcons.solidServer,
                                    color: COLOUR_PRIMARY,
                                  ),
                                ),
                                iconSize: 20,
                                elevation: 16,
                                style: TextStyle(
                                  fontSize: 16,
                                  color: COLOUR_PRIMARY,
                                ),
                                underline: Container(
                                  height: 1,
                                  color: COLOUR_PRIMARY,
                                ),
                                onChanged: (String newValue) {
                                  setState(() {
                                    hostDropDownItem = newValue;

                                    hostOption = hostDropDownItem;
                                    if (hostOption == "Custom") {
                                      showDialog(
                                        context: context,
                                        child: customHostDialog(),
                                      );
                                    } else {
                                      host = hostOptions[hostOption];

                                      updateGQLClient();
                                      versionCheck(context);
                                    }

                                    prefs.setString("host", host);
                                    prefs.setString("hostOption", hostOption);
                                  });
                                },
                                items: hostOptions.keys
                                    .map<DropdownMenuItem<String>>(
                                      (String value) =>
                                          DropdownMenuItem<String>(
                                        value: value,
                                        child: Text(value),
                                      ),
                                    )
                                    .toList(),
                              ),
                            ),
                          ],
                        ),

                        SizedBox(height: 14),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            loadingOverlay(loading),
          ],
        ),
      ),
    );
  }

  Widget forgotPasswordButton() {
    return Material(
      type: MaterialType.transparency,
      child: InkWell(
        child: Padding(
          padding: EdgeInsets.symmetric(vertical: 10),
          child: Text(
            "Forgot Password?",
            style: TextStyle(
              color: COLOUR_PRIMARY,
              fontWeight: FontWeight.w500,
              fontSize: 18,
            ),
          ),
        ),
        onTap: () {
          String email = emailInput.text.trim().toLowerCase();
          if (email.isEmpty || !isEmail(email)) {
            showWarningDialog(
              context,
              "",
              "Please enter a valid email address first.",
            );
          } else {
            showDialog(
              context: context,
              child: forgotPasswordDialog(),
            );
          }
        },
      ),
    );
  }

  Widget forgotPasswordDialog() {
    return SimpleDialog(
      title: Text("Forgot Password?"),
      contentPadding: EdgeInsets.all(15),
      children: <Widget>[
        textFieldWithIcon(
          "Email",
          keyboardType: TextInputType.emailAddress,
          controller: emailInput,
          validator: (value) {
            if (value.isEmpty)
              return ' This field is required ';
            else if (!isEmail(value))
              return ' Please enter a valid email address ';
            return null;
          },
          context: context,
          onFieldSubmitted: (_) => onForgotPassword(),
        ),
        SizedBox(height: 10),
        RaisedButton(
          child: Text(
            "Reset via SMS",
            style: TextStyle(
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          color: COLOUR_PRIMARY,
          onPressed: onForgotPassword,
        ),
      ],
    );
  }

  void setAPIHost() async {
    host = apiHostInput.text;
    updateGQLClient();
    Navigator.of(context).pop();
    versionCheck(context);
  }

  void setDefaultAPIHost() async {
    // Get Location
    Position position;
    try {
      position = await getGeolocatorPosition();
      if (position == null) {
        showErrorDialog(
          context,
          "Please give the app location service permission.",
        );
      }
    } catch (e) {
      return;
    }

    final coordinates = Coordinates(position.latitude, position.longitude);
    var addresses =
        await Geocoder.local.findAddressesFromCoordinates(coordinates);

    // Set host to China if not in Australia
    if (addresses.first.countryCode != "AU") {
      setState(() {
        hostOption = "China";
        host = hostOptions[hostOption];
        hostDropDownItem = hostOption;
      });
    }

    prefs.setString("host", host);
    prefs.setString("hostOption", hostOption);
  }

  Widget customHostDialog() {
    return SimpleDialog(
      title: Text("Custom API Host"),
      contentPadding: EdgeInsets.all(15),
      children: <Widget>[
        textFieldWithIcon(
          "Address",
          keyboardType: TextInputType.url,
          controller: apiHostInput,
          onFieldSubmitted: (_) => setAPIHost(),
        ),
        SizedBox(height: 10),
        RaisedButton(
          child: Text(
            "Set API Host",
            style: TextStyle(
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          color: COLOUR_PRIMARY,
          onPressed: setAPIHost,
        ),
      ],
    );
  }
}
