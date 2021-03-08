import 'package:awesome_dialog/awesome_dialog.dart';
import 'package:fieldapp/auth.dart';
import 'package:fieldapp/main.dart';
import 'package:fieldapp/pages/start.dart';
import 'package:fieldapp/scanner.dart';
import 'package:fieldapp/widgets/commonWidgets.dart';
import 'package:fieldapp/widgets/verifyBanner.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:package_info/package_info.dart';

import '../utils.dart';

/// Page Wrapper (Contains AppBar and Drawer)
class AppScaffold extends StatefulWidget {
  final Widget Function(
    bool Function() isLoading,
    void Function(bool) setLoading,
  ) builder;
  final String title;
  final double titleFontSize;
  AppScaffold({this.builder, this.title = "Genesis", this.titleFontSize = 25});

  @override
  AppScaffoldState createState() => new AppScaffoldState();
}

class AppScaffoldState extends State<AppScaffold> {
  String version = "";

  bool loading = false;
  bool isLoading() => loading;
  void setLoading(bool value) => setState(() => loading = value);

  bool hasScannerPlugin = true;
  bool continuousModeOn = false;

  @override
  void initState() {
    super.initState();
    SystemChrome.setSystemUIOverlayStyle(SystemUiOverlayStyle(
      statusBarIconBrightness: Brightness.light,
    ));

    getVersion();
    checkScannerPlugin();
  }

  void getVersion() async {
    PackageInfo packageInfo = await PackageInfo.fromPlatform();
    setState(() {
      // packageInfo.version is always 1.0.0 when debugging
      version =
          packageInfo.version == "1.0.0" ? "DEBUGGING" : packageInfo.version;
    });
  }

  void checkScannerPlugin() async {
    var response = await Scanner.getStatus();
    if (response.message == Scanner.MISSING_SCANNER_PLUGIN) {
      setState(() => hasScannerPlugin = false);
      return;
    }
    Scanner.setMethodCallHandler();
    Scanner.onContinuousSetCallback = onContinuousToggle;

    ScanMethodResponse r = await Scanner.getContinuousModeStatus();
    setState(() => continuousModeOn = r.value);
  }

  void onContinuousToggle(bool value) {
    setState(() => continuousModeOn = value);
  }

  Widget build(BuildContext context) {
    Widget page;
    Widget notification;

    if (continuousModeOn) {
      // Continuous mode on?
      notification = continuousModeIndicator();
    } else if (!me.mobileVerified &&
        me.mobilePhone != null &&
        me.mobilePhone.length > 0) {
      // Verify Email Message ?
      notification = VerifyBanner(
        message: "Please click here to verify your mobile",
        setLoading: setLoading,
      );
    }

    // Page
    if (notification != null) {
      page = Stack(
        children: [
          Padding(
            child: widget.builder(isLoading, setLoading),
            padding: EdgeInsets.only(top: 30),
          ),
          notification,
        ],
      );
    } else {
      page = widget.builder(isLoading, setLoading);
    }

    return Stack(
      children: <Widget>[
        startOverlayAnim(alignment: Alignment.topCenter, height: 78),
        Scaffold(
          backgroundColor: Colors.transparent,
          appBar: AppBar(
            backgroundColor: Colors.transparent,
            elevation: 0,
            title: Hero(
              tag: "title",
              child: Material(
                type: MaterialType.transparency,
                child: Text(
                  widget.title,
                  style: TextStyle(
                    fontSize: widget.titleFontSize,
                    color: Colors.white,
                    fontWeight: FontWeight.w600,
                  ),
                  maxLines: 1,
                ),
              ),
            ),
          ),
          drawer: Drawer(
            child: ListView(
              padding: EdgeInsets.zero,
              children: <Widget>[
                DrawerHeader(
                  decoration: BoxDecoration(
                    color: COLOUR_PRIMARY,
                  ),
                  padding: EdgeInsets.fromLTRB(16, 0, 16, 8),
                  child: Column(
                    children: <Widget>[
                      Container(
                        width: double.infinity,
                        child: Align(
                          child: CloseButton(
                            color: Colors.white,
                            onPressed: () {
                              Navigator.pop(context);
                            },
                          ),
                          alignment: Alignment.topRight,
                        ),
                      ),
                      Center(
                        child: Hero(
                          tag: "logo",
                          child: FaIcon(
                            FontAwesomeIcons.solidSteak,
                            color: Colors.white,
                            size: 55,
                          ),
                        ),
                      ),
                      SizedBox(height: 5),
                      Text(
                        'Genesis',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 15,
                          fontWeight: FontWeight.w600,
                          height: 0.8,
                        ),
                        maxLines: 1,
                      ),
                      Text(
                        version,
                        style: TextStyle(
                          fontSize: 9,
                          color: Colors.grey.shade700,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      Text(
                        "Server: " +
                            (hostOption == "Custom" ? host : hostOption),
                        style: TextStyle(
                          fontSize: 9,
                          color: Colors.grey.shade700,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ],
                  ),
                ),
                ListTile(
                  leading: FaIcon(
                    FontAwesomeIcons.usersClass,
                    color: COLOUR_PRIMARY,
                  ),
                  title: Text('Menu'),
                  onTap: () {
                    Navigator.pop(context);
                    navKey.currentState.pushNamed("/menu");
                  },
                ),
                ListTile(
                  leading: FaIcon(
                    FontAwesomeIcons.barcodeRead,
                    color: COLOUR_PRIMARY,
                  ),
                  title: Text('Scanner'),
                  onTap: () {
                    Navigator.pop(context);
                    navKey.currentState
                        .pushNamed(mode == null ? "/menu" : "/home");
                  },
                ),
                ListTile(
                  leading: FaIcon(
                    FontAwesomeIcons.userCircle,
                    color: COLOUR_PRIMARY,
                  ),
                  title: Text('Profile'),
                  onTap: () {
                    Navigator.pop(context);
                    navKey.currentState.pushNamed("/profile");
                  },
                ),
                ListTile(
                  leading: FaIcon(
                    FontAwesomeIcons.cog,
                    color: COLOUR_PRIMARY,
                  ),
                  title: Text('Settings'),
                  onTap: () {
                    Navigator.pop(context);
                    navKey.currentState.pushNamed("/settings");
                  },
                ),
                ListTile(
                  leading: FaIcon(
                    FontAwesomeIcons.wifi,
                    color: COLOUR_PRIMARY,
                  ),
                  title: Text('Check Internet'),
                  onTap: () async {
                    ConnectionStatus _conn = await checkConnectivity();
                    await AwesomeDialog(
                      context: context,
                      dialogType: _conn.isConnected
                          ? DialogType.SUCCES
                          : DialogType.ERROR,
                      tittle: "Check Internet",
                      desc: _conn.connectionMsg,
                      btnCancelOnPress: () {},
                      btnCancelColor: COLOUR_PRIMARY,
                      btnCancelText: "Ok",
                      headerAnimationLoop: false,
                    ).show();
                  },
                ),
                ListTile(
                  leading: FaIcon(
                    FontAwesomeIcons.signOut,
                    color: COLOUR_PRIMARY,
                  ),
                  title: Text('Log out'),
                  onTap: () async {
                    await AwesomeDialog(
                      context: context,
                      dialogType: DialogType.INFO,
                      tittle: "Are you sure?",
                      desc: "",
                      btnCancelOnPress: () {},
                      btnCancelColor: COLOUR_PRIMARY,
                      btnCancelText: "No",
                      btnOkText: "Log out",
                      btnOkOnPress: () {
                        mode = null;
                        Navigator.pop(context);
                        Auth.signOut();
                        navKey.currentState
                            .pushNamedAndRemoveUntil("/", (route) => false);
                      },
                      headerAnimationLoop: false,
                    ).show();
                  },
                ),
              ],
            ),
          ),
          body: GestureDetector(
            onTap: () => FocusScope.of(context)
                .requestFocus(FocusNode()), // close keyboard on tap outside
            child: Container(
              child: page,
              decoration: BoxDecoration(color: Colors.white),
              height: double.infinity,
              width: double.infinity,
            ),
          ),
        ),
        loadingOverlay(loading),
      ],
    );
  }

  Widget continuousModeIndicator() {
    return Stack(
      children: <Widget>[
        Container(
          height: 30,
          decoration: BoxDecoration(
            color: Colors.grey.shade600,
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.3),
                blurRadius: 3,
                offset: Offset(0, 2),
              )
            ],
          ),
        ),
        Material(
          type: MaterialType.transparency,
          child: InkWell(
            child: Container(
              height: 30,
              child: Center(
                child: Text(
                  "Continuous Mode On",
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }
}
