import 'package:awesome_dialog/awesome_dialog.dart';
import 'package:fieldapp/graphql/mutations.dart';
import 'package:fieldapp/graphql/queries.dart';
import 'package:fieldapp/main.dart';
import 'package:fieldapp/types/types.dart';
import 'package:fieldapp/utils.dart';
import 'package:fieldapp/widgets/dialogs.dart';
import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:geocoder/geocoder.dart';
import 'package:geohash/geohash.dart';
import 'package:geolocator/geolocator.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:smart_select/smart_select.dart';

import '../auth.dart';

TrackAction selectedTrackAction;
String memo;

void openTrackingMenuBasic({
  @required List<Product> products,
  @required List<Carton> cartons,
  @required List<Pallet> pallets,
  @required List<GenesisContainer> containers,
  @required BuildContext context,
  @required Function(bool) setLoading,
  @required @required String memo,
  String title,
}) async {
  if (me.role.trackActions.length == 0) {
    showErrorDialog(
      context,
      "You do not have permission to perform any track actions.",
    );
    return;
  }

  await AwesomeDialog(
    aligment: Alignment.topCenter,
    dismissOnTouchOutside: false,
    context: context,
    dialogType: DialogType.INFO,
    animType: AnimType.BOTTOMSLIDE,
    body: Container(
      transform: Matrix4.translationValues(0.0, -40.0, 0.0),
      child: TrackingMenu(
        dialogContext: context,
        title: title,
        onActionChange: (TrackAction action) => selectedTrackAction = action,
        onMemoChange: (String value) => memo = value,
      ),
    ),
    customHeader: FaIcon(
      FontAwesomeIcons.solidTruckMoving,
      size: 60,
    ),
  ).show();
}

class TrackMenuBasic extends StatefulWidget {
  @override
  _TrackMenuBasicState createState() => _TrackMenuBasicState();
}

class _TrackMenuBasicState extends State<TrackMenuBasic> {
  @override
  Widget build(BuildContext context) {
    return Container();
  }
}

// for non-advanced users
Future<bool> complete({
  @required List<Product> products,
  @required List<Carton> cartons,
  @required List<Pallet> pallets,
  @required List<GenesisContainer> containers,
  @required List<String> productScanTimes,
  @required List<String> cartonScanTimes,
  @required List<String> palletScanTimes,
  @required List<String> containerScanTimes,
  @required List<String> cartonPhotoBlobIDs,
  @required List<String> productPhotoBlobIDs,
  @required BuildContext context,
  @required Function(bool) setLoading,
}) async {
  if (selectedTrackAction == null) return false;

  setLoading(true);

  // Get Location
  Position position;
  try {
    position = await getGeolocatorPosition();

    if (position == null) {
      setLoading(false);
      showErrorDialog(
        context,
        "Please give the app location service permission.",
      );
    }
  } catch (e) {
    // Error?
    setLoading(false);
    print(e.toString());
    showErrorDialog(context, "Failed to get location");
    return false;
  }

  final coordinates = Coordinates(position.latitude, position.longitude);
  var addresses =
      await Geocoder.local.findAddressesFromCoordinates(coordinates);

  // Create location name
  String locationName = "";
  if (addresses.first.thoroughfare != null)
    locationName += addresses.first.thoroughfare;
  if (addresses.first.locality != null) {
    locationName = locationName.length > 0
        ? '$locationName, ${addresses.first.locality}'
        : addresses.first.locality;
  }
  if (addresses.first.adminArea != null) {
    locationName = locationName.length > 0
        ? '$locationName, ${addresses.first.adminArea}'
        : addresses.first.adminArea;
  }
  if (addresses.first.countryName != null) {
    locationName = locationName.length > 0
        ? '$locationName, ${addresses.first.countryName}'
        : addresses.first.countryName;
  }

  bool confirmed = false;
  await AwesomeDialog(
    context: context,
    dialogType: DialogType.INFO,
    animType: AnimType.BOTTOMSLIDE,
    tittle: "Confirm",
    body: WillPopScope(
      onWillPop: () async => false,
      child: Padding(
        padding: EdgeInsets.symmetric(horizontal: 10),
        child: Column(
          children: <Widget>[
            RichText(
              text: TextSpan(
                text:
                    "Commit the following Track Action to all products listed?",
                style: TextStyle(color: Colors.black),
                children: <TextSpan>[
                  TextSpan(
                    text:
                        " (including all products in any cartons, pallets or containers)",
                    style: TextStyle(
                      fontSize: 10,
                      color: Colors.grey.shade700,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(height: 20),
            Text(
              selectedTrackAction.name,
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
              textAlign: TextAlign.center,
            ),
            selectedTrackAction.nameChinese != ""
                ? Text(
                    selectedTrackAction.nameChinese,
                    style: TextStyle(color: Colors.grey.shade700),
                    textAlign: TextAlign.center,
                  )
                : Container(),
            Text(
              locationName,
              style: TextStyle(
                color: Colors.grey,
                fontSize: 12,
              ),
              textAlign: TextAlign.center,
            )
          ],
        ),
      ),
    ),
    btnCancelOnPress: () {
      setLoading(false);
      return false;
    },
    btnOkOnPress: () => confirmed = true,
    btnCancelColor: COLOUR_PRIMARY,
    headerAnimationLoop: false,
  ).show();

  if (!confirmed) {
    setLoading(false);
    return false;
  }

  // Mutation
  setLoading(true);

  try {
    bool timeout = false;
    QueryResult result = await client
        .mutate(
      MutationOptions(
        documentNode: gql(GQLMutation.recordTransaction),
        variables: {
          "input": {
            "trackActionCode": selectedTrackAction.code,
            "memo": memo,
            "productIDs": products.map((e) => e.id).toList(),
            "cartonIDs": cartons.map((e) => e.id).toList(),
            "palletIDs": pallets.map((e) => e.id).toList(),
            "containerIDs": containers.map((e) => e.id).toList(),
            "productScanTimes": productScanTimes,
            "cartonScanTimes": cartonScanTimes,
            "palletScanTimes": palletScanTimes,
            "containerScanTimes": containerScanTimes,
            "cartonPhotoBlobIDs": cartonPhotoBlobIDs,
            "productPhotoBlobIDs": productPhotoBlobIDs,
            "locationGeohash":
                Geohash.encode(position.latitude, position.longitude),
            "locationName": locationName,
          }
        },
      ),
    )
        .timeout(
      timeoutDuration,
      onTimeout: () {
        timeout = true;
        return QueryResult();
      },
    );

    setLoading(false);

    if (timeout) {
      showErrorDialog(context, "Timed out");
      return false;
    }

    // Error?
    if (result.hasException || result.data == null) {
      showErrorDialog(
        context,
        result.exception.graphqlErrors.length == 0
            ? "An issue occured"
            : result.exception.graphqlErrors[0].message,
      );
      return false;
    }
  } catch (e) {
    // Error?
    print(e.toString());
    showErrorDialog(context, "An issue occured");
  }

  setLoading(false);
  return true;
}

void openTrackingMenu({
  @required List<Product> products,
  @required List<Carton> cartons,
  @required List<Pallet> pallets,
  @required List<GenesisContainer> containers,
  @required List<String> productScanTimes,
  @required List<String> cartonScanTimes,
  @required List<String> palletScanTimes,
  @required List<String> containerScanTimes,
  @required List<String> cartonPhotoBlobIDs,
  @required List<String> productPhotoBlobIDs,
  @required BuildContext context,
  @required Function(bool) setLoading,
  @required Future<void> Function(dynamic) captureCartonImage,
  @required Future<void> Function(dynamic) captureProductImage,
}) async {
  if (me.role.trackActions.length == 0) {
    showErrorDialog(
      context,
      "You do not have permission to perform any track actions.",
    );
    return;
  }

  TrackAction selectedValue;
  String memo;

  await AwesomeDialog(
    aligment: Alignment.topCenter,
    context: context,
    dialogType: DialogType.INFO,
    animType: AnimType.BOTTOMSLIDE,
    body: Container(
      transform: Matrix4.translationValues(0.0, -40.0, 0.0),
      child: TrackingMenu(
        captureCartonImage: captureCartonImage,
        cartons: cartons,
        onActionChange: (TrackAction action) => selectedValue = action,
        onMemoChange: (String value) => memo = value,
      ),
    ),
    btnCancelOnPress: () => selectedValue = null,
    btnOkOnPress: () {},
    btnCancelColor: COLOUR_PRIMARY,
    customHeader: FaIcon(
      FontAwesomeIcons.solidTruckMoving,
      size: 60,
    ),
  ).show();

  if (selectedValue == null) return;

  setLoading(true);

  // Action require photo(s)?
  if (mode == Mode.advance &&
      selectedValue.requirePhotos.length > 1 &&
      (selectedValue.requirePhotos[0] || selectedValue.requirePhotos[1])) {
    if (cartons.length > 1) {
      // warning modal this action requires a photo + only  1 carton scanned
      AwesomeDialog(
        dismissOnTouchOutside: false,
        headerAnimationLoop: false,
        context: context,
        dialogType: DialogType.INFO,
        animType: AnimType.BOTTOMSLIDE,
        tittle: 'The track action that you have chosen requires only 1 carton',
        desc: "",
        btnOkOnPress: () async {
          setLoading(false);
          selectedTrackAction = null;
        },
      ).show();
      return;
    }

    if (cartons.length == 1) {
      // modal here to take a photo
      await AwesomeDialog(
        dismissOnTouchOutside: false,
        headerAnimationLoop: false,
        context: context,
        dialogType: DialogType.INFO,
        animType: AnimType.BOTTOMSLIDE,
        tittle:
            'The Track action you have chosen requires a photo of the carton',
        desc: "",
        btnOkText: "Take Photo",
        btnOkOnPress: () async {
          try {
            bool timeout = false;
            QueryResult result = await client
                .query(
              QueryOptions(
                documentNode: gql(GQLQuery.getObject),
                variables: {"id": cartons[0].id},
                fetchPolicy: FetchPolicy.networkOnly,
              ),
            )
                .timeout(
              timeoutDuration,
              onTimeout: () {
                timeout = true;
                return QueryResult();
              },
            );

            if (timeout) {
              showErrorDialog(context, "Timed out while getting object.");
              return;
            }

            // Error check
            if (result.hasException || result.data == null) {
              showErrorDialog(
                context,
                result.exception.graphqlErrors.length == 0
                    ? "An issue occurred."
                    : result.exception.graphqlErrors[0].message,
              );
              return;
            }

            // brings up camera
            if (selectedValue.requirePhotos[0]) {
              await captureCartonDialog(
                result.data["getObject"]["carton"],
                context,
                captureCartonImage,
                selectedTrackAction,
                setLoading,
              );
            } else {
              await captureProductDialog(
                result.data["getObject"]["carton"],
                context,
                captureProductImage,
                selectedTrackAction,
                setLoading,
              );
            }
          } catch (e) {
            Fluttertoast.showToast(
              msg: e.toString(),
              gravity: ToastGravity.TOP,
            );
            selectedTrackAction = null;
          }
        },
      ).show();
    }
  }

  // Get Location
  setLoading(true);
  Position position;
  try {
    position = await getGeolocatorPosition();
    if (position == null) {
      setLoading(false);
      showErrorDialog(
        context,
        "Please give the app location service permission.",
      );
    }
  } catch (e) {
    // Error?
    setLoading(false);
    print(e.toString());
    showErrorDialog(context, "Failed to get location");
    return;
  }

  final coordinates = Coordinates(position.latitude, position.longitude);
  var addresses =
      await Geocoder.local.findAddressesFromCoordinates(coordinates);

  String locationName =
      '${addresses.first.thoroughfare}, ${addresses.first.locality}, ${addresses.first.adminArea}, ${addresses.first.countryName}';

  // Confirm dialog
  bool confirmed = false;
  await AwesomeDialog(
    context: context,
    dialogType: DialogType.INFO,
    animType: AnimType.BOTTOMSLIDE,
    tittle: "Confirm",
    body: WillPopScope(
      onWillPop: () async => false,
      child: Padding(
        padding: EdgeInsets.symmetric(horizontal: 10),
        child: Column(
          children: <Widget>[
            RichText(
              text: TextSpan(
                text:
                    "Commit the following Track Action to all products listed?",
                style: TextStyle(color: Colors.black),
                children: <TextSpan>[
                  TextSpan(
                    text:
                        " (including all products in any cartons, pallets or containers)",
                    style: TextStyle(
                      fontSize: 10,
                      color: Colors.grey.shade700,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(height: 20),
            Text(
              selectedValue.name,
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
              textAlign: TextAlign.center,
            ),
            selectedValue.nameChinese != ""
                ? Text(
                    selectedValue.nameChinese,
                    style: TextStyle(color: Colors.grey.shade700),
                    textAlign: TextAlign.center,
                  )
                : Container(),
            Text(
              locationName,
              style: TextStyle(
                color: Colors.grey,
                fontSize: 12,
              ),
              textAlign: TextAlign.center,
            )
          ],
        ),
      ),
    ),
    btnCancelOnPress: () {
      setLoading(false);
      selectedTrackAction = null;
    },
    btnOkOnPress: () => confirmed = true,
    btnCancelColor: COLOUR_PRIMARY,
    headerAnimationLoop: false,
  ).show();

  if (!confirmed) return;

  // Mutation
  try {
    bool timeout = false;
    QueryResult result = await client
        .mutate(
      MutationOptions(
        documentNode: gql(GQLMutation.recordTransaction),
        variables: {
          "input": {
            "trackActionCode": selectedValue.code,
            "memo": memo,
            "productIDs": products.map((e) => e.id).toList(),
            "cartonIDs": cartons.map((e) => e.id).toList(),
            "palletIDs": pallets.map((e) => e.id).toList(),
            "containerIDs": containers.map((e) => e.id).toList(),
            "productScanTimes": productScanTimes,
            "cartonScanTimes": cartonScanTimes,
            "palletScanTimes": palletScanTimes,
            "containerScanTimes": containerScanTimes,
            "cartonPhotoBlobIDs": cartonPhotoBlobIDs,
            "productPhotoBlobIDs": productPhotoBlobIDs,
            "locationGeohash":
                Geohash.encode(position.latitude, position.longitude),
            "locationName": locationName,
          }
        },
      ),
    )
        .timeout(
      timeoutDuration,
      onTimeout: () {
        timeout = true;
        return QueryResult();
      },
    );

    if (timeout) {
      setLoading(false);
      showErrorDialog(context, "Timed out");
      return;
    }

    // Error?
    if (result.hasException || result.data == null) {
      setLoading(false);
      showErrorDialog(
        context,
        result.exception.graphqlErrors.length == 0
            ? "An issue occured"
            : result.exception.graphqlErrors[0].message,
      );

      selectedTrackAction = null;
      return;
    }

    // Success
    showSuccessDialog(
      context,
      "Success",
      "Tracking action successfully logged.",
    );
  } catch (e) {
    // Error?
    print(e.toString());
    showErrorDialog(context, "An issue occured");
  } finally {
    setLoading(false);
    selectedTrackAction = null;
  }
}

class TrackingMenu extends StatefulWidget {
  final Function(TrackAction) onActionChange;
  final Function(String) onMemoChange;
  final String title;
  @required
  final Future<void> Function(dynamic) captureCartonImage;
  @required
  final List<dynamic> cartons;
  final BuildContext dialogContext;

  TrackingMenu(
      {Key key,
      this.onActionChange,
      this.onMemoChange,
      this.title,
      this.cartons,
      this.captureCartonImage,
      this.dialogContext})
      : super(key: key);

  @override
  _TrackingMenuState createState() => _TrackingMenuState();
}

class _TrackingMenuState extends State<TrackingMenu> {
  TrackAction selection;

  List<SmartSelectOption<TrackAction>> options = me.role.trackActions
      .map((v) => SmartSelectOption<TrackAction>(value: v, title: v.name))
      .toList();

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: () async => false,
      child: SingleChildScrollView(
        child: Container(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: <Widget>[
              Container(
                child: Align(
                  child: CloseButton(
                      color: Colors.black,
                      onPressed: () {
                        Navigator.pop(context);
                        selectedTrackAction = null;
                        memo = null;
                      }),
                  alignment: Alignment.topRight,
                ),
              ),

              Center(
                child: Text(
                  widget.title != null ? widget.title : "Track Action",
                  style: TextStyle(
                    fontSize: 16,
                  ),
                ),
              ),

              // Track action
              SmartSelect<TrackAction>.single(
                title: "Track Action",
                value: selection,
                options: options,
                onChange: (value) {
                  setState(() => selection = value);
                  setState(() => selectedTrackAction = value);
                  widget.onActionChange(value);
                },
                builder: (
                  BuildContext context,
                  SmartSelectState<TrackAction> state,
                  SmartSelectShowModal showChoices,
                ) {
                  return FlatButton(
                    child: Padding(
                      padding: EdgeInsets.all(5),
                      child: state.value != null
                          ? trackActionCard(state.value)
                          : Text(
                              "Select an Action...",
                              style: TextStyle(
                                color: Colors.grey,
                                fontSize: 12,
                              ),
                            ),
                    ),
                    color: Colors.grey.shade200,
                    onPressed: () async {
                      await Auth.getMe();
                      showChoices(context);
                    },
                  );
                },
                choiceConfig: SmartSelectChoiceConfig<TrackAction>(
                  titleBuilder: (BuildContext context,
                      SmartSelectOption<TrackAction> item) {
                    return trackActionCard(item.value);
                  },
                ),
              ),
              SizedBox(height: 10),
              // Memo
              SingleChildScrollView(
                child: Padding(
                  padding: EdgeInsets.symmetric(horizontal: 10),
                  child: TextField(
                    onChanged: (value) {
                      setState(() => memo = value);
                      widget.onMemoChange(value);
                    },
                    decoration: InputDecoration(
                      labelText: "Memo (optional)",
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(0),
                      ),
                      hintStyle: TextStyle(fontSize: 14),
                    ),
                    maxLines: 2,
                    style: TextStyle(fontSize: 14),
                    textInputAction: TextInputAction.done,
                  ),
                ),
              ),
              mode != Mode.advance
                  ? Padding(
                      padding: EdgeInsets.symmetric(vertical: 20),
                      child: SizedBox(
                        width: 200,
                        height: 50,
                        child: (FlatButton(
                          disabledColor: Colors.grey.withOpacity(0.5),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(30.0),
                          ),
                          color: Color(0xFF00CA71),
                          child: Text(
                            "OK",
                            style: TextStyle(color: Colors.white),
                          ),
                          onPressed: selection == null
                              ? null
                              : () {
                                  Navigator.of(widget.dialogContext).pop();
                                  return;
                                },
                        )),
                      ),
                    )
                  : Container()
            ],
          ),
        ),
      ),
    );
  }

  Widget trackActionCard(TrackAction action) {
    if (action == null) return Container();
    if (action.nameChinese == "") return Text(action.name);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: <Widget>[
        Text(action.name),
        Text(
          action.nameChinese,
          style: TextStyle(
            color: Colors.grey,
            fontSize: 12,
          ),
        ),
      ],
    );
  }
}
