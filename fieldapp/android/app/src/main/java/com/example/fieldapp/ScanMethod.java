package com.genesis.fieldapp;

public enum ScanMethod {
    getStatus("getStatus"),
    turnOffScanner("turnOffScanner"),
    turnOnScanner("turnOnScanner"),
    resetScanner("resetScanner"),
  
    startScan("startScan"),
    stopScan("stopScan"),
  
    getBeepStatus("getBeepStatus"),
    setBeepOn("setBeepOn"),
    setBeepOff("setBeepOff"),
  
    getIndicatorLightStatus("getIndicatorLightStatus"),
    setIndicatorLightOn("setIndicatorLightOn"),
    setIndicatorLightOff("setIndicatorLightOff"),
  
    getVibrateStatus("getVibrateStatus"),
    setVibrateOn("setVibrateOn"),
    setVibrateOff("setVibrateOff"),

    getContinuousModeStatus("getContinuousModeStatus"),
    setContinuousModeOn("setContinuousModeOn"),
    setContinuousModeOff("setContinuousModeOff");


    private String str;
    ScanMethod(String s) { this.str = s; }
    public String toString() { return str; }
}