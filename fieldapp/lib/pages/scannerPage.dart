import 'dart:async';
import 'dart:io';
import 'package:flutter/services.dart';
import 'package:awesome_dialog/awesome_dialog.dart';
import 'package:fab_circular_menu/fab_circular_menu.dart';
import 'package:fieldapp/graphql/mutations.dart';
import 'package:fieldapp/graphql/queries.dart';
import 'package:fieldapp/main.dart';
import 'package:fieldapp/scanner.dart';
import 'package:fieldapp/types/types.dart';
import 'package:fieldapp/widgets/dialogs.dart';
import 'package:fieldapp/widgets/scannerBindMenus.dart';
import 'package:fieldapp/widgets/scannerTrackingMenu.dart';
import 'package:fieldapp/widgets/objectList.dart';
import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:qrscan/qrscan.dart' as cameraScanner;
import 'package:image_picker/image_picker.dart';
import 'package:http/http.dart' as http;
import 'package:http_parser/http_parser.dart';
import 'package:uuid/uuid.dart';
import '../utils.dart';

class ScannerPage extends StatefulWidget {
  final void Function(bool) setLoading;
  ScannerPage({Key key, @required this.setLoading}) : super(key: key);

  @override
  _ScannerPageState createState() => _ScannerPageState();
}

class _ScannerPageState extends State<ScannerPage> {
  List<String> uuids = List<String>();
  List<Product> products = List<Product>();
  List<Carton> cartons = List<Carton>();
  List<Pallet> pallets = List<Pallet>();
  List<GenesisContainer> containers = List<GenesisContainer>();
  List<String> productScanTimes = List<String>();
  List<String> cartonScanTimes = List<String>();
  List<String> palletScanTimes = List<String>();
  List<String> containerScanTimes = List<String>();
  List<String> cartonPhotoBlobIDs = List<String>();
  List<String> productPhotoBlobIDs = List<String>();

  String highlightUUID;
  Timer highlightTimer;

  FabCircularMenuController menuController = FabCircularMenuController();

  bool hasScannerPlugin = true;
  bool showScanButton = false;

  final imagePicker = ImagePicker();

  @override
  void initState() {
    super.initState();
    SystemChrome.setPreferredOrientations([DeviceOrientation.portraitUp]);

    if (prefs.containsKey("alwaysShowScanButton")) {
      setState(() {
        showScanButton = prefs.getBool("alwaysShowScanButton");
      });
    }
    checkScannerPlugin();
    clearLists(true);
    if (mode != null && mode != Mode.advance && selectedTrackAction == null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        trackingMenuBasic(null);
      });
    }
  }

  void checkScannerPlugin() async {
    var response = await Scanner.getStatus();
    if (response.message == Scanner.MISSING_SCANNER_PLUGIN) {
      setState(() => hasScannerPlugin = false);
      return;
    }

    Scanner.onScanCallback = onScanSKU;
  }

  // handles when user presses complete button on complete task mode
  void onCompleteTask() async {
    await AwesomeDialog(
            context: context,
            dialogType: DialogType.INFO,
            animType: AnimType.BOTTOMSLIDE,
            tittle: 'Confirm ${buildSuccessString()}',
            desc: "",
            btnOkOnPress: () {
              // confirm
              complete(
                context: context,
                products: products,
                cartons: cartons,
                pallets: pallets,
                containers: containers,
                productScanTimes: productScanTimes,
                cartonScanTimes: cartonScanTimes,
                palletScanTimes: palletScanTimes,
                containerScanTimes: containerScanTimes,
                cartonPhotoBlobIDs: cartonPhotoBlobIDs,
                productPhotoBlobIDs: productPhotoBlobIDs,
                setLoading: (bool value) => widget.setLoading(value),
              ).then((value) {
                // if failed/timeout or canceled
                if (value == false) return;

                AwesomeDialog(
                  context: context,
                  dialogType: DialogType.INFO,
                  animType: AnimType.BOTTOMSLIDE,
                  tittle: 'Succesfully ${buildSuccessString()}',
                  desc: "",
                  btnOkOnPress: () {},
                  btnCancelColor: COLOUR_PRIMARY,
                  btnCancelOnPress: () {},
                ).show();
                clearLists(true);
              });
            },
            btnCancelColor: COLOUR_PRIMARY,
            btnCancelOnPress: () {})
        .show();
  }

  void onPackCartonAndCompleteTask() async {
    // if user havent scanned any products or carton
    if (cartons.length <= 0 || products.length <= 0) {
      await AwesomeDialog(
        context: context,
        dialogType: DialogType.WARNING,
        animType: AnimType.BOTTOMSLIDE,
        tittle: "Please scan at least 1 product and 1 carton.",
        desc: "",
        btnOkOnPress: () {},
        btnCancelColor: COLOUR_PRIMARY,
      ).show();
      return;
    }

    if (cartons.length > 1) {
      await AwesomeDialog(
        context: context,
        dialogType: DialogType.WARNING,
        animType: AnimType.BOTTOMSLIDE,
        tittle: "Only 1 carton allowed for this mode.",
        desc: "",
        btnOkOnPress: () {},
        btnCancelColor: COLOUR_PRIMARY,
      ).show();
      return;
    }

    // confirm dialog
    await AwesomeDialog(
      context: context,
      dialogType: DialogType.INFO,
      animType: AnimType.BOTTOMSLIDE,
      tittle:
          'Confirm "${selectedTrackAction.name}" ${products.length} products into 1 carton.',
      desc: "",
      btnOkOnPress: () {
        // confirms track action
        complete(
          context: context,
          products: products,
          cartons: cartons,
          pallets: pallets,
          containers: containers,
          productScanTimes: productScanTimes,
          cartonScanTimes: cartonScanTimes,
          palletScanTimes: palletScanTimes,
          containerScanTimes: containerScanTimes,
          cartonPhotoBlobIDs: cartonPhotoBlobIDs,
          productPhotoBlobIDs: productPhotoBlobIDs,
          setLoading: (bool value) => widget.setLoading(value),
        ).then((value) {
          // if failed/timeout or canceled
          if (value == false) return;

          // binds product(s) to carton
          bindProductsToCarton(
            context: context,
            products: products,
            cartons: cartons,
            clearLists: () => clearLists(true),
            closeMenu: () => setState(() => menuController.isOpen = false),
            setLoading: (bool value) => widget.setLoading(value),
          );
          return;
        });
      },
      btnCancelColor: COLOUR_PRIMARY,
      btnCancelOnPress: () {
        return;
      },
    ).show();
  }

  bool isEmpty() =>
      products.isEmpty &&
      cartons.isEmpty &&
      pallets.isEmpty &&
      containers.isEmpty;

  Future<void> onScanSKU(String code) async => onScan(code, false);

  Future<void> onScan(String code, bool viaCamera) async {
    // force non advadced users to select a track action before scanning
    if (mode != Mode.advance && selectedTrackAction == null) {
      trackingMenuBasic(
        "Please select a track action.",
      );
      return;
    }
    if (code == null) return;
    // Get uuid
    String uuid;

    String splitBy = "/q/"; // General QR code (carton, pallet or container)

    // Product QR Code (view or register)?
    bool productQR = code.indexOf("?productID=") != -1;
    if (productQR) splitBy = "?productID=";

    List<String> codeSplit = code.split(splitBy);
    if (codeSplit.length != 2) {
      Fluttertoast.showToast(
        msg: "Invalid QR Code",
        gravity: ToastGravity.TOP,
      );
      return;
    }

    uuid = codeSplit[1];

    // register urls can have another argument (&wechat_id=) - trim it to get uuid
    int aIndex = uuid.indexOf("&");
    if (productQR && aIndex != -1) uuid = uuid.substring(0, aIndex);

    // already scanned uuid
    if (uuids.contains(uuid)) {
      if (viaCamera) {
        showScanNextDialog(context, "This QR has already been scanned.", "",
            () async {
          String result = await cameraScanner.scan();
          onScan(result, viaCamera);
        });
      }
      setState(() => highlightUUID = uuid);
      if (highlightTimer != null) highlightTimer.cancel();
      highlightTimer = Timer(const Duration(milliseconds: 500), () {
        setState(() {
          highlightUUID = null;
          highlightTimer = null;
        });
      });
      return;
    }

    widget.setLoading(true);

    try {
      // query object
      bool timeout = false;
      QueryResult result = await client
          .query(
        QueryOptions(
          documentNode: gql(GQLQuery.getObject),
          variables: {"id": uuid},
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

      widget.setLoading(false);

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

      // only scan cartons or products in "pack carton and complete mode"
      if (mode == Mode.packCartonCompleteTask) {
        // check if scanned item neither product or carton
        if (result.data["getObject"]["carton"] == null &&
            result.data["getObject"]["product"] == null) {
          AwesomeDialog(
            headerAnimationLoop: false,
            context: context,
            dialogType: DialogType.WARNING,
            animType: AnimType.BOTTOMSLIDE,
            tittle: "Please scan a product or a carton.",
            desc: "",
            btnOkOnPress: () {},
            btnOkColor: COLOUR_PRIMARY,
          ).show();
          return;
        }
      }

      // check if carton
      if (result.data["getObject"]["carton"] != null) {
        // check if already bound to a pallet
        if (result.data["getObject"]["carton"]["pallet"] != null) {
          await showInfoDialog(context,
              "Cannot scan, item already bound to a pallet", "", () {});
          return;
        }

        // opens image dialog/ modal for user to take photo after scan
        if (selectedTrackAction != null &&
            selectedTrackAction.requirePhotos.length > 1) {
          if (selectedTrackAction.requirePhotos[0]) {
            captureCartonDialog(
              result.data["getObject"]["carton"],
              context,
              captureCartonImage,
              selectedTrackAction,
              widget.setLoading,
            );
          } else if (selectedTrackAction.requirePhotos[1]) {
            captureProductDialog(
              result.data["getObject"]["carton"],
              context,
              captureProductImage,
              selectedTrackAction,
              widget.setLoading,
            );
          }
        }
      }

      if (await addObject(result.data["getObject"]))
        setState(() => uuids.add(uuid));
    } catch (e) {
      // Error?
      showErrorDialog(context, "An issue occurred: " + e.toString());
    } finally {
      widget.setLoading(false);
    }
  }

  // uploads an image(blob) to db and returns the blob id
  Future<GQLResponse> uploadImage() async {
    var pickedImage = await imagePicker.getImage(
      source: ImageSource.camera,
      maxWidth: 1280,
      imageQuality: 85,
    );

    // response of image upload
    GQLResponse response = new GQLResponse();

    // get file
    if (pickedImage == null) {
      response.message = "No photo taken";
      return response;
    }
    var image = File(pickedImage.path);

    // byte data from captured image
    var byteData = image.readAsBytesSync();

    // image to multipart file
    var multipartFile = http.MultipartFile.fromBytes(
      'photo',
      byteData,
      filename: '${Uuid().v4().toString()}.jpg',
      contentType: MediaType("image", "jpg"),
    );

    // check image is null
    if (image != null) {
      widget.setLoading(true);

      try {
        bool timeout = false;
        QueryResult result = await client
            .mutate(
          MutationOptions(
            documentNode: gql(GQLMutation.fileUpload),
            variables: {
              'file': multipartFile,
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
          response.message = "Timed out";
          return response;
        }

        // Error check
        if (result.hasException || result.data == null) {
          if (result.exception.graphqlErrors.length == 0) return response;
          response.message = result.exception.graphqlErrors[0].toString();
          return response;
        }

        response.success = true;
        response.message = result.data["fileUpload"]["id"];
      } catch (e) {
        print(e.toString());
      } finally {
        widget.setLoading(false);
      }
    }
    return response;
  }

  Future<void> captureCartonImage(dynamic carton) async {
    // capture/upload image
    GQLResponse uploadResponse = await uploadImage();

    if (!uploadResponse.success) {
      showErrorDialog(context, uploadResponse.message);

      // remove carton from list if failed
      int index = cartons.indexWhere((v) => v.id == carton['id']);
      setState(() {
        cartons.removeAt(index);
        cartonScanTimes.removeAt(index);
        uuids.remove(carton['id']);
      });
      return;
    }

    setState(() => cartonPhotoBlobIDs.add(uploadResponse.message));

    await showSuccessDialog(
      context,
      "Success",
      "Carton photo uploaded.",
    ).then((value) {
      captureProductDialog(
        carton,
        context,
        captureProductImage,
        selectedTrackAction,
        (bool value) => widget.setLoading(value),
      );
    });
  }

  Future<void> captureProductImage(dynamic carton) async {
    // capture/upload image
    GQLResponse uploadResponse = await uploadImage();

    if (!uploadResponse.success) {
      showErrorDialog(context, uploadResponse.message);

      // remove carton from list if failed
      int index = cartons.indexWhere((v) => v.id == carton['id']);
      setState(() {
        cartons.removeAt(index);
        cartonScanTimes.removeAt(index);
        uuids.remove(carton['id']);
      });
      return;
    }

    setState(() => productPhotoBlobIDs.add(uploadResponse.message));

    // Success
    showSuccessDialog(
      context,
      "Success",
      "Product photo uploaded.",
    );
  }

  Future<bool> addObject(dynamic object) async {
    // Set object
    if (object["product"] != null) {
      setState(() {
        products.add(Product.fromJson(object["product"]));
        productScanTimes.add(DateTime.now().toUtc().toIso8601String());
      });
      return true;
    }
    if (object["carton"] != null) {
      setState(() {
        cartons.add(Carton.fromJson(object["carton"]));
        cartonScanTimes.add(DateTime.now().toUtc().toIso8601String());
      });
      return true;
    }
    if (object["pallet"] != null) {
      setState(() {
        pallets.add(Pallet.fromJson(object["pallet"]));
        palletScanTimes.add(DateTime.now().toUtc().toIso8601String());
      });
      return true;
    }
    if (object["container"] != null) {
      setState(() {
        containers.add(GenesisContainer.fromJson(object["container"]));
        containerScanTimes.add(DateTime.now().toUtc().toIso8601String());
      });
      return true;
    }

    return false;
  }

  Future<void> refreshList() async {
    if (isEmpty()) return;
    try {
      bool timeout = false;
      QueryResult result = await client
          .query(
        QueryOptions(
          documentNode: gql(GQLQuery.getObjects),
          variables: {
            'input': {
              'productIDs': products.map((e) => e.id).toList(),
              'cartonIDs': cartons.map((e) => e.id).toList(),
              'palletIDs': pallets.map((e) => e.id).toList(),
              'containerIDs': containers.map((e) => e.id).toList(),
            },
          },
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
        showErrorDialog(context, "Timed out while refreshing list.");
        return;
      }

      // Error check
      if (result.hasException ||
          result.data == null ||
          result.data["getObjects"] == null) {
        showErrorDialog(
          context,
          result.exception.graphqlErrors.length == 0
              ? "An issue occurred when trying to refresh."
              : result.exception.graphqlErrors[0].message,
        );
        return;
      }

      // Set object
      if (result.data["getObjects"]["products"] != null) {
        setState(
          () => products =
              (result.data["getObjects"]["products"] as List<dynamic>)
                  .map((e) => Product.fromJson(e))
                  .toList(),
        );
      }
      if (result.data["getObjects"]["cartons"] != null) {
        setState(
          () => cartons =
              (result.data["getObjects"]["cartons"] as List<dynamic>)
                  .map((e) => Carton.fromJson(e))
                  .toList(),
        );
      }
      if (result.data["getObjects"]["pallets"] != null) {
        setState(
          () => pallets =
              (result.data["getObjects"]["pallets"] as List<dynamic>)
                  .map((e) => Pallet.fromJson(e))
                  .toList(),
        );
      }
      if (result.data["getObjects"]["containers"] != null) {
        setState(
          () => containers =
              (result.data["getObjects"]["containers"] as List<dynamic>)
                  .map((e) => GenesisContainer.fromJson(e))
                  .toList(),
        );
      }
    } catch (e) {
      showErrorDialog(context, "An issue occurred when trying to refresh.");
    }
  }

  // brings up dialog to clear list
  void clearListsDialog() {
    AwesomeDialog(
      context: context,
      dialogType: DialogType.INFO,
      animType: AnimType.BOTTOMSLIDE,
      tittle: 'Clear List?',
      desc: "",
      btnCancelOnPress: () {},
      btnOkOnPress: () => clearLists(false),
      btnCancelColor: COLOUR_PRIMARY,
    ).show();
  }

  // clears item lists
  void clearLists(bool clearMemoTrackAction) {
    setState(() {
      uuids.clear();
      products.clear();
      cartons.clear();
      pallets.clear();
      containers.clear();
      productScanTimes.clear();
      cartonScanTimes.clear();
      palletScanTimes.clear();
      containerScanTimes.clear();
      cartonPhotoBlobIDs.clear();
      productPhotoBlobIDs.clear();
      if (clearMemoTrackAction) {
        selectedTrackAction = null;
        memo = null;
      }
    });
  }

  // modified tracking menu for non-advanced users
  void trackingMenuBasic(String title) {
    openTrackingMenuBasic(
      memo: memo,
      title: title,
      context: context,
      products: products,
      cartons: cartons,
      pallets: pallets,
      containers: containers,
      setLoading: (bool value) => widget.setLoading(value),
    );
  }

  String buildSuccessString() {
    String s = """ "${selectedTrackAction.name}" to:
                 ${containers.length == 0 ? "" : containers.length.toString() + " container(s), "}
                 ${pallets.length == 0 ? "" : pallets.length.toString() + " pallet(s), "}
                 ${cartons.length == 0 ? "" : cartons.length.toString() + " carton(s), "}
                 ${products.length == 0 ? "" : products.length.toString() + " product(s), "}""";

    return stripMargin(s.substring(0, s.length - 2));
  }

  @override
  Widget build(BuildContext context) {
    bool empty = isEmpty();
    Color iconColour = empty ? Colors.white.withOpacity(0.4) : Colors.white;
    return WillPopScope(
      onWillPop: () async => false,
      child: Scaffold(
        body: mode == Mode.advance
            ? // show wheel menu if advanced mode
            FabCircularMenu(
                controller: menuController,
                ringColor: COLOUR_SECONDARY.withOpacity(0.9),
                ringDiameter: MediaQuery.of(context).size.width * 0.8,
                child: Stack(
                  children: <Widget>[
                    pleaseScanMessage(),
                    RefreshIndicator(
                      onRefresh: refreshList,
                      child: SingleChildScrollView(
                        child: Column(
                          children: <Widget>[
                            ObjectList(
                              header: "Container",
                              list: containers,
                              dismissItem: (i) => setState(() {
                                uuids.remove(containers[i].id);
                                containers.removeAt(i);
                              }),
                              highlightUUID: highlightUUID,
                            ),
                            ObjectList(
                              header: "Pallet",
                              list: pallets,
                              dismissItem: (i) => setState(() {
                                uuids.remove(pallets[i].id);
                                pallets.removeAt(i);
                              }),
                              highlightUUID: highlightUUID,
                            ),
                            ObjectList(
                              header: "Carton",
                              list: cartons,
                              dismissItem: (i) => setState(() {
                                uuids.remove(cartons[i].id);
                                cartons.removeAt(i);
                              }),
                              highlightUUID: highlightUUID,
                            ),
                            ObjectList(
                              header: "Product",
                              list: products,
                              dismissItem: (i) => setState(() {
                                uuids.remove(products[i].id);
                                products.removeAt(i);
                              }),
                              highlightUUID: highlightUUID,
                            ),
                            SizedBox(height: 90),
                          ],
                        ),
                      ),
                    ),
                    hasScannerPlugin && !showScanButton
                        ? Container()
                        : Align(
                            alignment: Alignment.bottomLeft,
                            child: Padding(
                              padding: EdgeInsets.only(bottom: 30, left: 15),
                              child: FloatingActionButton.extended(
                                heroTag: "scanQRCodeButton",
                                icon: FaIcon(FontAwesomeIcons.barcodeScan,
                                    color: Colors.white),
                                label: Text('Scan QR Code'),
                                tooltip: 'Scan QR Code',
                                onPressed: () async {
                                  String result = await cameraScanner.scan();
                                  onScan(result, true);
                                },
                              ),
                            ),
                          ),
                  ],
                ),
                options: <Widget>[
                  IconButton(
                    icon: FaIcon(
                      FontAwesomeIcons.solidBroom,
                      color: iconColour,
                    ),
                    tooltip: "Clear List",
                    onPressed: empty ? null : clearListsDialog,
                  ),
                  IconButton(
                    icon: FaIcon(
                      FontAwesomeIcons.solidUnlink,
                      color: iconColour,
                    ),
                    tooltip: "Unbind",
                    onPressed: () => empty
                        ? null
                        : openUnbindMenu(
                            context: context,
                            products: products,
                            cartons: cartons,
                            pallets: pallets,
                            containers: containers,
                            closeMenu: () =>
                                setState(() => menuController.isOpen = false),
                            setLoading: (bool value) =>
                                widget.setLoading(value),
                          ),
                  ),
                  IconButton(
                    icon: FaIcon(
                      FontAwesomeIcons.solidLink,
                      color: iconColour,
                    ),
                    tooltip: "Bind",
                    onPressed: () => empty
                        ? null
                        : openBindMenu(
                            context: context,
                            products: products,
                            cartons: cartons,
                            pallets: pallets,
                            containers: containers,
                            closeMenu: () =>
                                setState(() => menuController.isOpen = false),
                            setLoading: (bool value) =>
                                widget.setLoading(value),
                            activateScan: () async {
                              String result = await cameraScanner.scan();
                              if (result == null) return;
                              return await onScan(result, true);
                            },
                          ),
                  ),
                ],
              )
            : // for non-advaced users
            Container(
                child: Stack(
                  children: <Widget>[
                    pleaseScanMessage(),
                    RefreshIndicator(
                      onRefresh: refreshList,
                      child: SingleChildScrollView(
                        child: Column(
                          children: <Widget>[
                            ObjectList(
                              header: "Container",
                              list: containers,
                              dismissItem: (i) => setState(() {
                                uuids.remove(containers[i].id);
                                containers.removeAt(i);
                              }),
                              highlightUUID: highlightUUID,
                            ),
                            ObjectList(
                              header: "Pallet",
                              list: pallets,
                              dismissItem: (i) => setState(() {
                                uuids.remove(pallets[i].id);
                                pallets.removeAt(i);
                              }),
                              highlightUUID: highlightUUID,
                            ),
                            ObjectList(
                              header: "Carton",
                              list: cartons,
                              dismissItem: (i) => setState(() {
                                uuids.remove(cartons[i].id);
                                cartons.removeAt(i);
                              }),
                              highlightUUID: highlightUUID,
                            ),
                            ObjectList(
                              header: "Product",
                              list: products,
                              dismissItem: (i) => setState(() {
                                uuids.remove(products[i].id);
                                products.removeAt(i);
                              }),
                              highlightUUID: highlightUUID,
                            ),
                            SizedBox(height: 90),
                          ],
                        ),
                      ),
                    ),
                    Row(
                      children: <Widget>[
                        hasScannerPlugin && !showScanButton
                            ? Container()
                            : Align(
                                alignment: Alignment.bottomLeft,
                                child: Padding(
                                  padding:
                                      EdgeInsets.only(bottom: 15, left: 10),
                                  child: FloatingActionButton.extended(
                                    heroTag: "scanQRCodeButton",
                                    icon: FaIcon(
                                      FontAwesomeIcons.barcodeScan,
                                      color: Colors.white,
                                    ),
                                    label: Text(
                                      MediaQuery.of(context).size.width <= 360
                                          ? 'Scan QR'
                                          : 'Scan QR Code',
                                    ),
                                    tooltip: 'Scan QR Code',
                                    onPressed: () async {
                                      // force non advanced users to select a track action before scanning if they havent already
                                      if (mode != Mode.advance &&
                                          selectedTrackAction == null) {
                                        trackingMenuBasic(
                                          "Please select a track action.",
                                        );
                                        return;
                                      }

                                      String result =
                                          await cameraScanner.scan();
                                      onScan(result, true);
                                    },
                                  ),
                                ),
                              ),
                        Align(
                          alignment: Alignment.bottomCenter,
                          child: Padding(
                            padding: EdgeInsets.only(bottom: 15, left: 10),
                            child: FloatingActionButton.extended(
                              backgroundColor:
                                  isEmpty() ? Colors.grey : Colors.black,
                              heroTag: "completeButton",
                              icon: FaIcon(
                                FontAwesomeIcons.check,
                                color: Colors.white,
                              ),
                              label: Text('Complete'),
                              tooltip: 'Complete',
                              onPressed: () async {
                                // if user has not scanned any items
                                if (isEmpty()) return;

                                // pack carton and complete task mode
                                if (mode == Mode.packCartonCompleteTask) {
                                  onPackCartonAndCompleteTask();
                                }

                                // complete task mode
                                if (mode == Mode.completeTask) {
                                  onCompleteTask();
                                }
                              },
                            ),
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
      ),
    );
  }

  Widget pleaseScanMessage() {
    return Container(
      width: double.infinity,
      child: AnimatedOpacity(
        opacity: isEmpty() ? 0.5 : 0,
        duration: Duration(milliseconds: 250),
        child: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              SizedBox(height: 60),
              Container(
                padding: EdgeInsets.symmetric(horizontal: 20),
                child: FaIcon(
                  FontAwesomeIcons.solidQrcode,
                  size: 200,
                ),
                decoration: BoxDecoration(
                  border: Border.fromBorderSide(
                    BorderSide(
                      color: COLOUR_PRIMARY,
                      width: 8,
                    ),
                  ),
                ),
              ),
              SizedBox(height: 15),
              Text(
                "Please Scan a QR Code",
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 24,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
