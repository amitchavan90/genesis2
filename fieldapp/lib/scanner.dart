import 'package:flutter/services.dart';

enum ScanMethod {
  getStatus,
  turnOffScanner,
  turnOnScanner,
  resetScanner,

  startScan,
  stopScan,

  getBeepStatus,
  setBeepOn,
  setBeepOff,

  getIndicatorLightStatus,
  setIndicatorLightOn,
  setIndicatorLightOff,

  getVibrateStatus,
  setVibrateOn,
  setVibrateOff,

  getContinuousModeStatus,
  setContinuousModeOn,
  setContinuousModeOff,
}

class ScanMethodResponse {
  bool value = false;
  bool success = false;
  String message = "";
}

class Scanner {
  static const channelSKUMethod = const MethodChannel('genesis.sku/method');

  static const String MISSING_SCANNER_PLUGIN = "Missing Scanner Plugin.";

  static Function(String) onScanCallback;
  static Function(bool) onContinuousSetCallback;

  /// Invoke a Java Scanner SDK method
  static Future<ScanMethodResponse> invokeMethod(ScanMethod method,
      {bool boolean = false}) async {
    ScanMethodResponse response = new ScanMethodResponse();

    try {
      String methodString = method.toString().replaceAll("ScanMethod.", "");

      if (boolean) {
        final bool result = await channelSKUMethod.invokeMethod(methodString);
        response.success = true;
        response.value = result;
      } else {
        final String result = await channelSKUMethod.invokeMethod(methodString);
        response.success = true;
        response.message = result;
      }
    } on PlatformException catch (e) {
      response.success = false;
      response.message = "Failed to Invoke: '${e.message}'.";
    } on MissingPluginException {
      response.success = false;
      response.message = MISSING_SCANNER_PLUGIN;
    }

    return response;
  }

  static void setMethodCallHandler() {
    channelSKUMethod.setMethodCallHandler((MethodCall call) {
      switch (call.method) {
        case "ON_RECEIVE":
          if (onScanCallback == null || call.arguments == null) return;
          onScanCallback(call.arguments.toString());
          break;

        case "CONTINUOUS_SET":
          if (onContinuousSetCallback == null || call.arguments == null) return;
          onContinuousSetCallback(call.arguments as bool);
          break;

        default:
          print("unknown callback: ${call.method}");
      }

      return;
    });
  }

  static Future<ScanMethodResponse> getStatus() async =>
      invokeMethod(ScanMethod.getStatus, boolean: true);
  static Future<ScanMethodResponse> turnOn() async =>
      invokeMethod(ScanMethod.turnOnScanner);
  static Future<ScanMethodResponse> turnOff() async =>
      invokeMethod(ScanMethod.turnOffScanner);
  static Future<ScanMethodResponse> reset() async =>
      invokeMethod(ScanMethod.resetScanner);

  static Future<ScanMethodResponse> startScan() async =>
      invokeMethod(ScanMethod.startScan);
  static Future<ScanMethodResponse> stopScan() async =>
      invokeMethod(ScanMethod.stopScan);

  static Future<ScanMethodResponse> getBeepState() async =>
      invokeMethod(ScanMethod.getBeepStatus, boolean: true);
  static Future<ScanMethodResponse> setBeepOn() async =>
      invokeMethod(ScanMethod.setBeepOn);
  static Future<ScanMethodResponse> setBeepOff() async =>
      invokeMethod(ScanMethod.setBeepOff);

  static Future<ScanMethodResponse> getIndicatorLightStatus() async =>
      invokeMethod(ScanMethod.getIndicatorLightStatus, boolean: true);
  static Future<ScanMethodResponse> setIndicatorLightOn() async =>
      invokeMethod(ScanMethod.setIndicatorLightOn);
  static Future<ScanMethodResponse> setIndicatorLightOff() async =>
      invokeMethod(ScanMethod.setIndicatorLightOff);

  static Future<ScanMethodResponse> getVibrateStatus() async =>
      invokeMethod(ScanMethod.getVibrateStatus, boolean: true);
  static Future<ScanMethodResponse> setVibrateOn() async =>
      invokeMethod(ScanMethod.setVibrateOn);
  static Future<ScanMethodResponse> setVibrateOff() async =>
      invokeMethod(ScanMethod.setVibrateOff);

  static Future<ScanMethodResponse> getContinuousModeStatus() async =>
      invokeMethod(ScanMethod.getContinuousModeStatus, boolean: true);
  static Future<ScanMethodResponse> setContinuousModeOn() async =>
      invokeMethod(ScanMethod.setContinuousModeOn);
  static Future<ScanMethodResponse> setContinuousModeOff() async =>
      invokeMethod(ScanMethod.setContinuousModeOff);
}
