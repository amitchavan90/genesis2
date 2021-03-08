import 'dart:developer';

import 'package:fieldapp/auth.dart';
import 'package:fieldapp/types/types.dart';
import 'package:flutter/painting.dart';
import 'package:flutter/widgets.dart';
import 'package:fieldapp/widgets/dialogs.dart';
import 'package:fieldapp/main.dart';
import 'package:flutter/material.dart';
import 'package:package_info/package_info.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:ota_update/ota_update.dart';

void versionCheck(BuildContext context) async {
  PackageInfo packageInfo = await PackageInfo.fromPlatform();

  // packageInfo.version is always 1.0.0 when debugging
  if (packageInfo.version == "1.0.0") return;

  GQLResponse response = await Auth.appVersionCheck();

  // Error check
  if (!response.success) {
    log(response.toString());
    log(response.message);
    log(response.hashCode.toString());

    showErrorDialog(context, "An issue occured when checking server version.");
    return;
  }

  // Version match
  if (packageInfo.version == response.message) return;

  // Check version
  final serverVersions = response.message.split(".");
  final localVersions = packageInfo.version.split(".");
  if (serverVersions.length == localVersions.length) {
    var needsUpdate = false;
    for (var i = 0; i < serverVersions.length; i++) {
      var s = int.tryParse(serverVersions[i]);
      var l = int.tryParse(localVersions[i]);
      if (s > l) {
        needsUpdate = true;
        break;
      }
    }
    if (!needsUpdate) return;
  }

  // Show update dialog
  var fileName = 'genesis-${response.message}.apk';
  var downloadLink = 'https://admin.gn.latitude28.cn/fieldapp/$fileName';

  showDialog(
    context: context,
    child: SimpleDialog(
      title: Text("Update Required"),
      contentPadding: EdgeInsets.all(15),
      children: <Widget>[
        Text("Your app is out of date."),
        Text("Latest: ${response.message}"),
        Text("Current: ${packageInfo.version}"),
        SizedBox(height: 10),
        RaisedButton(
          child: Text(
            "Update",
            style: TextStyle(
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          color: COLOUR_PRIMARY,
          onPressed: () {
            Navigator.of(context).pop();

            showDialog(
              context: context,
              barrierDismissible: false,
              child: SimpleDialog(
                title: Text("Updating"),
                contentPadding: EdgeInsets.all(15),
                children: <Widget>[
                  UpdaterDialog(downloadLink, fileName),
                ],
              ),
            );
          },
        ),
        FlatButton(
          child: Text(fileName),
          onPressed: () async {
            var url = downloadLink;
            if (await canLaunch(url)) {
              await launch(url);
            } else {
              throw 'Could not launch $url';
            }
          },
        ),
      ],
    ),
    barrierDismissible: false,
  );
}

class UpdaterDialog extends StatefulWidget {
  final String downloadLink;
  final String fileName;
  UpdaterDialog(this.downloadLink, this.fileName);

  @override
  _UpdaterDialogState createState() => _UpdaterDialogState();
}

class _UpdaterDialogState extends State<UpdaterDialog> {
  OtaEvent otaEvent;

  @override
  void initState() {
    super.initState();

    try {
      OtaUpdate()
          .execute(
        widget.downloadLink,
        destinationFilename: widget.fileName,
      )
          .listen(
        (OtaEvent event) {
          if (event.status != OtaStatus.DOWNLOADING) {
            Navigator.of(context).pop();
            return;
          }

          setState(() => otaEvent = event);
        },
      );
    } catch (e) {
      showErrorDialog(
        context,
        'Failed to make OTA update. Details: $e',
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    if (otaEvent == null || otaEvent.status != OtaStatus.DOWNLOADING)
      return Center(
        child: Column(
          children: [
            otaEvent != null
                ? Text(otaEvent.status.toString().replaceFirst("OtaEvent.", ""))
                : Container(),
            CircularProgressIndicator()
          ],
        ),
      );

    var progress = int.tryParse(otaEvent.value);

    return Center(
      child: Column(
        children: [
          progress != null
              ? LinearProgressIndicator(value: progress / 100.0)
              : Container(),
        ],
      ),
    );
  }
}
