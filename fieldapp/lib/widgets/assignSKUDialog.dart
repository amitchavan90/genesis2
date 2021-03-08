import 'package:awesome_dialog/awesome_dialog.dart';
import 'package:fieldapp/main.dart';
import 'package:fieldapp/types/enums.dart';
import 'package:fieldapp/types/types.dart';
import 'package:fieldapp/widgets/dialogs.dart';
import 'package:fieldapp/widgets/skuSelect.dart';
import 'package:flutter/material.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:fieldapp/graphql/mutations.dart';

/// Prompts for SKU selection then converts blanks to products, returning the new products
void assignSKUDialog(
  BuildContext context, {
  @required List<Product> products,
  @required void setLoading(bool value),
}) async {
  if (products.length == 0) return;

  // Get SKU
  SKU selectedValue;
  await AwesomeDialog(
    context: context,
    dialogType: DialogType.INFO,
    animType: AnimType.BOTTOMSLIDE,
    body: SKUSelect(onChange: (value) => selectedValue = value),
    btnCancelOnPress: () => selectedValue = null,
    btnOkOnPress: () {},
    btnCancelColor: COLOUR_PRIMARY,
    headerAnimationLoop: false,
  ).show();

  if (selectedValue == null) return;

  String plural = products.length > 1 ? 's' : '';

  // Mutation
  setLoading(true);
  try {
    bool timeout = false;
    QueryResult result = await client
        .mutate(
      MutationOptions(
        documentNode: gql(GQLMutation.productBatchAction),
        variables: {
          'ids': products.map((e) => e.id).toList(),
          'action':
              BatchAction.SetSKU.toString().replaceFirst('BatchAction.', ''),
          'value': {'str': selectedValue.id},
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

      return;
    }

    // Success
    products.forEach((p) => p.sku = selectedValue);

    showSuccessDialog(
      context,
      "Success",
      'Successfully set product$plural SKU to: ${selectedValue.code}',
    );
  } catch (e) {
    // Error?
    print(e.toString());
    showErrorDialog(context, "An issue occured");
  }
  setLoading(false);
  return;
}
