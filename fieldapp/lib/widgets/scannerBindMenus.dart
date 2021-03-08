import 'package:awesome_dialog/awesome_dialog.dart';
import 'package:fieldapp/graphql/mutations.dart';
import 'package:fieldapp/main.dart';
import 'package:fieldapp/types/enums.dart';
import 'package:fieldapp/types/types.dart';
import 'package:fieldapp/widgets/assignSKUDialog.dart';
import 'package:fieldapp/widgets/dialogs.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:smart_select/smart_select.dart';

void openBindMenu({
  @required BuildContext context,
  @required List<Product> products,
  @required List<Carton> cartons,
  @required List<Pallet> pallets,
  @required List<GenesisContainer> containers,
  @required Function closeMenu,
  @required Function(bool) setLoading,
  @required Future<void> Function() activateScan,
}) {
  openListDialog(
    context: context,
    title: "Bind",
    icon: FontAwesomeIcons.solidLink,
    children: <Widget>[
      _bindOption(
        context,
        products,
        cartons,
        FontAwesomeIcons.lightSteak,
        FontAwesomeIcons.box,
        title: 'Product${products.length == 1 ? '' : 's'} to Carton',
        onPressed: () => _doAction(
          GQLMutation.productBatchAction,
          products,
          BatchAction.SetCarton,
          value: cartons,
          itemName: "Product",
          otherItemName: "Carton",
          successMessage: "attached to",
          context: context,
          closeMenu: closeMenu,
          setLoading: setLoading,
        ),
        otherItemName: "Carton",
        activateScan: activateScan,
      ),
      _bindOption(
        context,
        cartons,
        pallets,
        FontAwesomeIcons.box,
        FontAwesomeIcons.palletAlt,
        title: 'Carton${cartons.length == 1 ? '' : 's'} to Pallet',
        onPressed: () => _doAction(
          GQLMutation.cartonBatchAction,
          cartons,
          BatchAction.SetPallet,
          value: pallets,
          itemName: "Carton",
          otherItemName: "Pallet",
          successMessage: "attached to",
          context: context,
          closeMenu: closeMenu,
          setLoading: setLoading,
        ),
        otherItemName: "Pallet",
        activateScan: activateScan,
        color: Colors.grey.shade100,
      ),
      _bindOption(
        context,
        pallets,
        containers,
        FontAwesomeIcons.palletAlt,
        FontAwesomeIcons.containerStorage,
        title: 'Pallet${cartons.length == 1 ? '' : 's'} to Container',
        onPressed: () => _doAction(
          GQLMutation.palletBatchAction,
          pallets,
          BatchAction.SetContainer,
          value: containers,
          itemName: "Pallet",
          otherItemName: "Container",
          successMessage: "attached to",
          context: context,
          closeMenu: closeMenu,
          setLoading: setLoading,
        ),
        otherItemName: "Container",
        activateScan: activateScan,
      ),
      _bindOption(
        context,
        products,
        null,
        FontAwesomeIcons.barcodeAlt,
        FontAwesomeIcons.lightSteak,
        title: 'SKU to Product${products.length == 1 ? '' : 's'}',
        onPressed: () {
          closeMenu();
          assignSKUDialog(
            context,
            products: products,
            setLoading: setLoading,
          );
        },
        needListB: false,
      )
    ],
  );
}

void bindProductsToCarton({
  @required BuildContext context,
  @required List<Product> products,
  @required List<Carton> cartons,
  @required Function closeMenu,
  @required Function(bool) setLoading,
  @required Function() clearLists,
}) {
  _bindProductsToCartons(
      GQLMutation.productBatchAction, products, BatchAction.SetCarton,
      value: cartons,
      itemName: "Product",
      successMessage: "attached to",
      context: context,
      closeMenu: closeMenu,
      setLoading: setLoading,
      clearLists: clearLists);
}

void openUnbindMenu({
  @required BuildContext context,
  @required List<Product> products,
  @required List<Carton> cartons,
  @required List<Pallet> pallets,
  @required List<GenesisContainer> containers,
  @required Function closeMenu,
  @required Function(bool) setLoading,
}) {
  openListDialog(
    context: context,
    title: "Unbind",
    icon: FontAwesomeIcons.solidUnlink,
    children: <Widget>[
      _bindOption(
        context,
        products,
        cartons,
        FontAwesomeIcons.lightSteak,
        FontAwesomeIcons.box,
        title: 'Product${products.length == 1 ? '' : 's'} from Carton(s)',
        onPressed: () => _doAction(
          GQLMutation.productBatchAction,
          products,
          BatchAction.DetachFromCarton,
          itemName: "Product",
          otherItemName: "Carton",
          successMessage: "detached from their Carton(s).",
          context: context,
          closeMenu: closeMenu,
          setLoading: setLoading,
        ),
        needListB: false,
      ),
      _bindOption(
        context,
        cartons,
        pallets,
        FontAwesomeIcons.box,
        FontAwesomeIcons.palletAlt,
        title: 'Carton${cartons.length == 1 ? '' : 's'} from Pallet(s)',
        onPressed: () => _doAction(
          GQLMutation.cartonBatchAction,
          cartons,
          BatchAction.DetachFromPallet,
          itemName: "Carton",
          otherItemName: "Pallet",
          successMessage: "detached from their Pallet(s).",
          context: context,
          closeMenu: closeMenu,
          setLoading: setLoading,
        ),
        color: Colors.grey.shade100,
        needListB: false,
      ),
      _bindOption(
        context,
        pallets,
        containers,
        FontAwesomeIcons.palletAlt,
        FontAwesomeIcons.containerStorage,
        title: 'Pallet${cartons.length == 1 ? '' : 's'} from Container(s)',
        onPressed: () => _doAction(
          GQLMutation.palletBatchAction,
          pallets,
          BatchAction.DetachFromContainer,
          itemName: "Pallet",
          otherItemName: "Container",
          successMessage: "detached from their Container(s).",
          context: context,
          closeMenu: closeMenu,
          setLoading: setLoading,
        ),
        needListB: false,
      )
    ],
  );
}

void _doAction(
  String mutation,
  List<GenesisObject> list,
  BatchAction action, {
  List<GenesisObject> value,
  String successMessage,
  String itemName,
  String otherItemName,
  @required BuildContext context,
  @required Function closeMenu,
  @required Function(bool) setLoading,
}) async {
  if (list.length == 0) return;

  closeMenu();

  // Get value
  GenesisObject selectedValue;

  if (value != null && value.length > 1) {
    selectedValue = value[0];
    // Multiple values? (eg: binding products to carton but have multiple cartons scanned in? ask user which one)
    List<SmartSelectOption<GenesisObject>> options = value
        .map((v) => SmartSelectOption<GenesisObject>(value: v, title: v.code))
        .toList();

    await AwesomeDialog(
      context: context,
      dialogType: DialogType.INFO,
      animType: AnimType.BOTTOMSLIDE,
      body: StatefulBuilder(
        builder: (ctx, _setState) {
          GenesisObject _selection = selectedValue;

          return Column(
            children: <Widget>[
              Text('Which $otherItemName?'),
              SmartSelect<GenesisObject>.single(
                title: otherItemName,
                value: _selection,
                options: options,
                onChange: (val) {
                  _setState(() => _selection = val);
                  selectedValue = val;
                },
              ),
            ],
          );
        },
      ),
      btnCancelOnPress: () => selectedValue = null,
      btnOkOnPress: () {},
      btnCancelColor: COLOUR_PRIMARY,
    ).show();
    if (selectedValue == null) return;
  } else if (value != null && value.length == 1) {
    // Bind Confirm
    await AwesomeDialog(
      context: context,
      dialogType: DialogType.INFO,
      animType: AnimType.BOTTOMSLIDE,
      tittle: 'Confirm',
      desc: 'Bind $itemName${list.length == 1 ? '' : 's'} to ${value[0].code}?',
      btnCancelOnPress: () {},
      btnOkOnPress: () => selectedValue = value[0],
      btnCancelColor: COLOUR_PRIMARY,
    ).show();

    if (selectedValue == null) return;
  } else if (value == null) {
    // Unbind Confirm
    bool confirmed = false;
    await AwesomeDialog(
      context: context,
      dialogType: DialogType.INFO,
      animType: AnimType.BOTTOMSLIDE,
      tittle: 'Confirm',
      desc:
          'Unbind $itemName${list.length == 1 ? '' : 's'} from $otherItemName(s)?',
      btnCancelOnPress: () {},
      btnOkOnPress: () => confirmed = true,
      btnCancelColor: COLOUR_PRIMARY,
    ).show();
    if (!confirmed) return;
  } else {
    showWarningDialog(
      context,
      "No ${otherItemName}s scanned",
      "Please scan the $otherItemName you want to bind to.",
    );
    return;
  }

  // Mutation
  setLoading(true);
  try {
    Map<String, dynamic> variables = {
      'ids': list.map((e) => e.id).toList(),
      'action': action.toString().replaceFirst('BatchAction.', ''),
    };
    if (selectedValue != null) variables['value'] = {'str': selectedValue.id};

    bool timeout = false;
    QueryResult result = await client
        .mutate(
      MutationOptions(
        documentNode: gql(mutation),
        variables: variables,
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

      return;
    }

    // Success
    if (list[0] is Product) {
      for (Product obj in list) obj.carton = selectedValue;
    } else if (list[0] is Carton) {
      for (Carton obj in list) obj.pallet = selectedValue;
    } else if (list[0] is Pallet) {
      for (Pallet obj in list) obj.container = selectedValue;
    }

    if (successMessage != null)
      showSuccessDialog(
        context,
        "Success",
        '$itemName${list.length == 1 ? '' : 's'} successfully $successMessage' +
            (selectedValue != null ? ' ${selectedValue.code}' : ''),
      );
  } catch (e) {
    // Error?
    print(e.toString());
    showErrorDialog(context, "An issue occured");
  }
  setLoading(false);
}

void _bindProductsToCartons(
  String mutation,
  List<GenesisObject> list,
  BatchAction action, {
  List<GenesisObject> value,
  String successMessage,
  String itemName,
  @required BuildContext context,
  @required Function closeMenu,
  @required Function(bool) setLoading,
  @required Function() clearLists,
}) async {
  if (list.length == 0) return;
  closeMenu();

  // Get value
  GenesisObject selectedValue;

  selectedValue = value[0];
  if (selectedValue == null) return;

  // Mutation
  setLoading(true);

  try {
    Map<String, dynamic> variables = {
      'ids': list.map((e) => e.id).toList(),
      'action': action.toString().replaceFirst('BatchAction.', ''),
    };
    if (selectedValue != null) variables['value'] = {'str': selectedValue.id};

    bool timeout = false;
    QueryResult result = await client
        .mutate(
      MutationOptions(
        documentNode: gql(mutation),
        variables: variables,
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

      return;
    }

    if (list[0] is Product) {
      for (Product obj in list) obj.carton = selectedValue;
    } else if (list[0] is Carton) {
      for (Carton obj in list) obj.pallet = selectedValue;
    }

    if (successMessage != null)
      // success
      showSuccessDialog(
        context,
        "Success",
        '$itemName${list.length == 1 ? '' : 's'} successfully $successMessage' +
            (selectedValue != null ? ' carton (${selectedValue.code})' : ''),
      );
  } catch (e) {
    // Error?
    print(e.toString());
    showErrorDialog(context, "An issue occured");
  }
  setLoading(false);
  clearLists();
}

Widget _bindOption(
  context,
  List<GenesisObject> listA,
  List<GenesisObject> listB,
  IconData iconA,
  IconData iconB, {
  String title,
  String subtitle,
  Function onPressed,
  Color color,
  bool needListB = true,
  String otherItemName,
  Future<void> Function() activateScan,
}) {
  return FlatButton(
    onPressed: listA.length > 0
        ? () {
            Navigator.of(context).pop(); // close pop-up

            if (needListB && listB.length == 0) {
              AwesomeDialog(
                context: context,
                dialogType: DialogType.INFO,
                animType: AnimType.BOTTOMSLIDE,
                body: Text("Please scan $otherItemName QR Code"),
                btnCancelOnPress: () {},
                btnOkOnPress: () async {
                  await activateScan();
                  if (onPressed != null) onPressed();
                },
                btnCancelColor: COLOUR_PRIMARY,
                customHeader: FaIcon(
                  iconB,
                  size: 60,
                ),
                btnOkText: "Scan",
              ).show();

              return;
            }

            if (onPressed != null) onPressed();
          }
        : null,
    color: Colors.white,
    child: Padding(
      child: Row(
        children: <Widget>[
          FaIcon(iconA),
          subtitle != null
              ? Column(
                  children: <Widget>[
                    Text(
                      title,
                      style: TextStyle(fontSize: 14),
                    ),
                    Text(
                      subtitle,
                      style: TextStyle(
                        fontSize: 9,
                        color: COLOUR_PRIMARY,
                      ),
                    )
                  ],
                )
              : Text(
                  title,
                  style: TextStyle(fontSize: 14),
                ),
          FaIcon(iconB),
        ],
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
      ),
      padding: EdgeInsets.symmetric(vertical: 12),
    ),
  );
}
