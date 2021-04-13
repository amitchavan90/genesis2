import 'dart:async';
import 'dart:io';
import 'dart:convert' as convert;
import 'package:connectivity/connectivity.dart';
import 'package:fieldapp/main.dart';
import 'package:flutter/widgets.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:path_provider/path_provider.dart';
import 'package:http/http.dart' as http;
import 'package:geolocator/geolocator.dart';

bool isEmail(String value) => value.length > 2 && value.indexOf("@") > 0;

String formatNumber(String number) {
  RegExp reg = new RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))');
  Function mathFunc = (Match match) => '${match[1]},';
  return number.replaceAllMapped(reg, mathFunc);
}

Color getColourFromHex(String hexColor) {
  hexColor = hexColor.toUpperCase().replaceAll("#", "");
  if (hexColor.length == 6) hexColor = "FF" + hexColor;
  return Color(int.parse(hexColor, radix: 16));
}

List<Color> hexListToColourList(List<dynamic> list) {
  List<Color> colourList = [];
  list.forEach((dynamic str) {
    colourList.add(getColourFromHex(str as String));
  });
  return colourList;
}

List<String> colourListToHexList(List<Color> list) {
  List<String> hexList = [];
  list.forEach((Color colour) {
    hexList.add(colour.toString().substring(8, 16).toUpperCase());
  });
  return hexList;
}

List<int> stringToByteArray(String str) {
  try {
    return convert.base64.decode(str);
  } catch (Exception) {}

  return null;
}

String byteArrayToString(List<int> ary) {
  try {
    return convert.base64.encode(ary);
  } catch (Exception) {
    return null;
  }
}

Future<File> downloadFile(String fileName, List<int> bytes) async {
  Directory tempDir = await getTemporaryDirectory();
  String tempPath = tempDir.path;
  File file = File('$tempPath/$fileName');
  return await file.writeAsBytes(bytes);
}

class StringToEnum {
  static String parse(enumItem) {
    if (enumItem == null) return null;
    return enumItem.toString().split('.')[1];
  }

  static T fromString<T>(List<T> enumValues, String value) {
    if (value == null || enumValues == null) return null;

    return enumValues.singleWhere(
        (enumItem) => parse(enumItem)?.toLowerCase() == value?.toLowerCase(),
        orElse: () => null);
  }

  static List<T> fromList<T>(List<T> enumValues, List<String> values) {
    if (values == null || values.length == 0 || enumValues == null) return null;

    List<T> result = new List<T>();
    for (int i = 0; i < values.length; i++) {
      result.add(
        enumValues.singleWhere(
            (enumItem) =>
                parse(enumItem)?.toLowerCase() == values[i]?.toLowerCase(),
            orElse: () => null),
      );
    }

    return result;
  }

  static List<String> toList<T>(List<T> enumValues) {
    if (enumValues == null) return null;
    var enumList = enumValues.map((t) => parse(t)).toList();
    return enumList;
  }
}

String stripMargin(String s) {
  return s.splitMapJoin(
    RegExp(r'^', multiLine: true),
    onMatch: (_) => ' ',
    onNonMatch: (n) => n.trim(),
  );
}

// ConnectionStatus iniformation of wifi/internet connection
class ConnectionStatus {
  String connectionMsg;
  bool isConnected;

  ConnectionStatus(String connectionMsg, bool isConnected) {
    this.connectionMsg = connectionMsg;
    this.isConnected = isConnected;
  }
}

Future<ConnectionStatus> checkConnectivity() async {
  ConnectionStatus conn = ConnectionStatus("", false);

  var connect = await (Connectivity().checkConnectivity());
  if (connect == ConnectivityResult.mobile) {
    conn.connectionMsg = "connected";
    conn.isConnected = true;
  }
  if (connect == ConnectivityResult.wifi) {
    conn.connectionMsg = "connected";
    conn.isConnected = true;
  }
  if (connect == ConnectivityResult.none) {
    conn.connectionMsg = "no internet connection or internet turned off";
    conn.isConnected = false;
    return conn;
  }

  try {
    final result = await http.get(host);

    if (result.statusCode == 200) {
      conn.connectionMsg = "connected";
      conn.isConnected = true;
    }

    if (result.statusCode != 200) {
      conn.connectionMsg = "no internet connection";
      conn.isConnected = false;
      return conn;
    }
  } on SocketException catch (_) {
    conn.connectionMsg = "no internet connection (server down)";
    conn.isConnected = false;
  } on ClientException catch (_) {
    conn.connectionMsg = "no internet connection";
    conn.isConnected = false;
  }
  return conn;
}

/// Gets the device's current gps location (while checking/asking for location service permission)
Future<Position> getGeolocatorPosition() async {
  // Get Location Permission
  LocationPermission perm = await Geolocator.checkPermission();
  if (perm == LocationPermission.denied) {
    perm = await Geolocator.requestPermission();
  }
  if (perm == LocationPermission.denied ||
      perm == LocationPermission.deniedForever) {
    return null;
  }

  // Get location
  return await Geolocator.getCurrentPosition(
    desiredAccuracy: LocationAccuracy.high,
    timeLimit: Duration(seconds: 20),
  );
}
