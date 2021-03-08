import 'package:fieldapp/graphql/queries.dart';
import 'package:fieldapp/main.dart';
import 'package:fieldapp/types/types.dart';
import 'package:fieldapp/widgets/dialogs.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:smart_select/smart_select.dart';

class SKUSelect extends StatefulWidget {
  final void Function(SKU value) onChange;
  SKUSelect({Key key, this.onChange}) : super(key: key);

  @override
  _SKUSelectState createState() => _SKUSelectState();
}

class _SKUSelectState extends State<SKUSelect> {
  List<SKU> skus;
  bool loading = true;
  SKU selection;

  @override
  void initState() {
    super.initState();

    getList(null);
  }

  Future<List<SmartSelectOption<SKU>>> getList(String search) async {
    bool timeout = false;
    QueryResult result = await client
        .query(
      QueryOptions(
        documentNode: gql(GQLQuery.skus),
        variables: {
          'search': {
            'search': search,
            'filter': 'Active',
            'sortDir': 'Descending'
          },
          'limit': 15,
          'offset': 0,
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
      showErrorDialog(context, "Timed out when trying to fetch SKUs.");
      return null;
    }

    // Error check
    if (result.hasException ||
        result.data == null ||
        result.data["skus"] == null) {
      showErrorDialog(
        context,
        result.exception.graphqlErrors.length == 0
            ? "An issue occurred when trying to fetch SKUs."
            : result.exception.graphqlErrors[0].message,
      );
      return null;
    }

    setState(() {
      skus = (result.data["skus"]["skus"] as List<dynamic>)
          .map((e) => SKU.fromJson(e))
          .toList();
      loading = false;
    });

    return SmartSelectOption.listFrom<SKU, SKU>(
      source: skus,
      value: (index, item) => item,
      title: (index, item) => item.name,
      subtitle: (index, item) => item.code,
    );
  }

  Widget build(BuildContext context) {
    return Column(
      children: <Widget>[Text('Which SKU?'), selectList()],
    );
  }

  Widget selectList() {
    if (loading) return CircularProgressIndicator();
    if (skus == null) return Text('Failed to fetch SKUs');
    return SmartSelect<SKU>.single(
      title: "SKU",
      value: selection,
      options: SmartSelectOption.listFrom<SKU, SKU>(
        source: skus,
        value: (index, item) => item,
        title: (index, item) => item.name,
        subtitle: (index, item) => item.code,
      ),
      onChange: (val) {
        setState(() => selection = val);
        if (widget.onChange != null) widget.onChange(val);
      },
      onSearch: (value) => getList(value),
      modalConfig: SmartSelectModalConfig(
        useFilter: true,
      ),
      builder: (
        BuildContext context,
        SmartSelectState<SKU> state,
        SmartSelectShowModal showChoices,
      ) {
        return FlatButton(
          child: Padding(
            padding: EdgeInsets.all(5),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: <Widget>[
                Text("SKU"),
                state.value != null
                    ? skuCard(state.value)
                    : Text(
                        "Select a SKU...",
                        style: TextStyle(
                          color: Colors.grey,
                          fontSize: 12,
                        ),
                      ),
              ],
            ),
          ),
          color: Colors.grey.shade200,
          onPressed: () => showChoices(context),
        );
      },
      choiceConfig: SmartSelectChoiceConfig<SKU>(
        titleBuilder: (BuildContext context, SmartSelectOption<SKU> item) {
          return skuCard(item.value);
        },
        subtitleBuilder: (BuildContext context, SmartSelectOption<SKU> item) {
          return null;
        },
      ),
    );
  }
}

Widget skuCard(SKU sku) {
  if (sku == null) return Container();
  return Row(
    children: <Widget>[
      SizedBox(
        width: 80,
        child: Center(
          child: sku.masterPlanURL != null
              ? Image.network(
                  host + sku.masterPlanURL,
                  height: 50,
                )
              : FaIcon(sku.icon),
        ),
      ),
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: <Widget>[
          Text(sku.name),
          Text(
            sku.code,
            style: TextStyle(
              color: Colors.grey,
              fontSize: 12,
            ),
          ),
        ],
      )
    ],
  );
}
