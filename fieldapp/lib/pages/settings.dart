import 'package:fieldapp/main.dart';
import 'package:flutter/material.dart';
import 'package:fieldapp/scanner.dart';

class SettingsPage extends StatefulWidget {
  SettingsPage({Key key}) : super(key: key);

  @override
  _SettingsPageState createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  bool beepStatus = false;
  bool vibrateStatus = false;
  bool indicatorLightStatus = false;
  bool alwaysShowScanButton = false;

  bool hasScannerPlugin = true;

  @override
  void initState() {
    super.initState();

    if (prefs.containsKey("alwaysShowScanButton"))
      setState(() {
        alwaysShowScanButton = prefs.getBool("alwaysShowScanButton");
      });

    loadSettings();
  }

  void loadSettings() async {
    // Load settings
    ScanMethodResponse resp = await Scanner.getBeepState();

    if (!resp.success) {
      // Failed to get scanner setting - ignore rest
      setState(() => hasScannerPlugin = false);
      return;
    }
    setState(() => beepStatus = resp.value);

    Scanner.getVibrateStatus()
        .then((r) => setState(() => vibrateStatus = r.value));
    Scanner.getIndicatorLightStatus()
        .then((r) => setState(() => indicatorLightStatus = r.value));
  }

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: EdgeInsets.all(20),
      children: <Widget>[
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: <Widget>[
            Text(
              "Scanner Options",
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
              ),
            ),
            hasScannerPlugin
                ? Container()
                : Text(
                    "(no scanner plugin detected)",
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.grey,
                    ),
                  ),
          ],
        ),
        Container(
          decoration: BoxDecoration(
            border: Border.all(
              color: COLOUR_PRIMARY.withOpacity(0.5),
              width: 2,
            ),
            borderRadius: BorderRadius.circular(5),
            color: hasScannerPlugin ? null : COLOUR_PRIMARY.withOpacity(0.2),
          ),
          padding: EdgeInsets.symmetric(horizontal: 10),
          child: Column(
            children: <Widget>[
              switchItem(
                "Beep on Scan",
                beepStatus,
                (value) {
                  setState(() => beepStatus = value);
                  if (value)
                    Scanner.setBeepOn();
                  else
                    Scanner.setBeepOff();
                },
                disabled: !hasScannerPlugin,
              ),
              switchItem(
                "Vibrate on Scan",
                vibrateStatus,
                (value) {
                  setState(() => vibrateStatus = value);
                  if (value)
                    Scanner.setVibrateOn();
                  else
                    Scanner.setVibrateOff();
                },
                disabled: !hasScannerPlugin,
              ),
              switchItem(
                "Indicator Light on Scan",
                indicatorLightStatus,
                (value) {
                  setState(() => indicatorLightStatus = value);
                  if (value)
                    Scanner.setIndicatorLightOn();
                  else
                    Scanner.setIndicatorLightOff();
                },
                disabled: !hasScannerPlugin,
              ),
              SizedBox(height: 5),
              switchItem(
                'Always show "Scan QR Code" Button',
                alwaysShowScanButton,
                (value) {
                  setState(() => alwaysShowScanButton = value);
                  prefs.setBool("alwaysShowScanButton", value);
                },
                disabled: !hasScannerPlugin,
              ),
              Text(
                "(will show scan button and allow scanning via phone camera even with scanner plugin present)",
                style: TextStyle(
                  fontSize: 10,
                  color: Colors.grey,
                ),
              ),
              SizedBox(height: 10),
            ],
          ),
        ),
      ],
    );
  }

  Widget switchItem(String name, bool value, Function(bool) onChanged,
      {bool disabled = false}) {
    return Row(
      children: <Widget>[
        Expanded(
          child: Text(
            name,
            style: TextStyle(
              fontSize: 18,
              color: disabled ? Colors.grey : null,
            ),
          ),
        ),
        Switch(
          value: value,
          onChanged: disabled
              ? null
              : (v) {
                  if (v == value) return;
                  onChanged(v);
                },
        ),
      ],
    );
  }
}
