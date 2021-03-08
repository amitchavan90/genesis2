package com.genesis.fieldapp;

import androidx.annotation.NonNull;
import io.flutter.embedding.android.FlutterActivity;
import io.flutter.embedding.engine.FlutterEngine;
import io.flutter.plugin.common.MethodChannel;

import android.device.ScanDevice;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.view.KeyEvent;

public class MainActivity extends FlutterActivity {
    private static final String CHANNEL_SKU_METHOD = "genesis.sku/method";
    private static final String SCAN_ACTION = "scan.rcv.message";
    private static final String METHOD_ON_RECEIVE = "ON_RECEIVE";
    private static final String METHOD_CONTINUOUS_SET = "CONTINUOUS_SET";

    ScanDevice sm;
    MethodChannel channelSKUMethod;

    @Override
    public void configureFlutterEngine(@NonNull FlutterEngine flutterEngine) {
        super.configureFlutterEngine(flutterEngine);
        try {
            sm = new ScanDevice();
        } catch (Exception e) {
            return;
        }
        sm.setOutScanMode(0); // broadcast mode
        
        channelSKUMethod = new MethodChannel(flutterEngine.getDartExecutor().getBinaryMessenger(), CHANNEL_SKU_METHOD);
        channelSKUMethod.setMethodCallHandler((call, result) -> {
            boolean success = false;

            switch (ScanMethod.valueOf(call.method)) {

                case getStatus:
                    result.success(sm.isScanOpened());
                    return;
                case turnOffScanner:
                    success = sm.closeScan();
                    break;
                case turnOnScanner:
                    success = sm.openScan();
                    break;
                case resetScanner:
                    success = sm.resetScan();
                    break;

                case startScan:
                    success = sm.startScan();
                    break;
                case stopScan:
                    success = sm.stopScan();
                    break;

                case getBeepStatus:
                    result.success(sm.getScanBeepState());
                    return;
                case setBeepOn:
                    success = sm.setScanBeep();
                    break;
                case setBeepOff:
                    success = sm.setScanUnBeep();
                    break;

                case getIndicatorLightStatus:
                    result.success(sm.getIndicatorLightMode() == 1);
                    return;
                case setIndicatorLightOn:
                    sm.setIndicatorLightMode(1);
                    success = true;
                    break;
                case setIndicatorLightOff:
                    sm.setIndicatorLightMode(0);
                    success = true;
                    break;

                case getVibrateStatus:
                    result.success(sm.getScanVibrateState());
                    return;
                case setVibrateOn:
                    success = sm.setScanVibrate();
                    break;
                case setVibrateOff:
                    success = sm.setScanUnVibrate();
                    break;

                case getContinuousModeStatus:
                    result.success(sm.getScanLaserMode() == 4);
                    return;
                case setContinuousModeOn:
                    sm.setScanLaserMode(4);
                    success = true;
                    break;
                case setContinuousModeOff:
                    sm.setScanLaserMode(8);
                    success = true;
                    break;
            }

            result.success(call.method + (success ? " Succeeded" : "Failed"));

        });
    }

    private BroadcastReceiver mScanReceiver = new BroadcastReceiver() {
        @Override
        public void onReceive(Context context, Intent intent) {
            if (intent.getAction().equalsIgnoreCase((SCAN_ACTION))) {
                byte[] barocode = intent.getByteArrayExtra("barocode");
                int barcodeLen = intent.getIntExtra("length", 0);
                // byte barcodeType = intent.getByteExtra("barcodeType", (byte) 0);
                // byte[] aimid = intent.getByteArrayExtra("aimid");
                String barcodeStr = new String(barocode, 0, barcodeLen);

                channelSKUMethod.invokeMethod(METHOD_ON_RECEIVE, barcodeStr);
                sm.stopScan();
            }
        }
    };

    @Override
    public boolean onKeyUp(int keyCode, KeyEvent event) {
        if (keyCode == KeyCode.SCAN) {
            // Toggle Continuous Scan
            int value = sm.getScanLaserMode() == 4 ? 8 : 4;
            sm.setScanLaserMode(value);
            channelSKUMethod.invokeMethod(METHOD_CONTINUOUS_SET, value == 4);
        } else if (keyCode == KeyCode.SCAN_LEFT || keyCode == KeyCode.SCAN_RIGHT) {
            // Turn-off Continuous Scan
            if (sm.getScanLaserMode() == 4) {
                sm.setScanLaserMode(8);
                channelSKUMethod.invokeMethod(METHOD_CONTINUOUS_SET, false);
            }
        }
        return super.onKeyUp(keyCode, event);
    }

    @Override
    protected void onPause() {
        super.onPause();
        if (sm != null) {
            sm.stopScan();
        }
        unregisterReceiver(mScanReceiver);
    }

    @Override
    protected void onResume() {
        super.onResume();
        IntentFilter filter = new IntentFilter();
        filter.addAction(SCAN_ACTION);
        registerReceiver(mScanReceiver, filter);
    }
}