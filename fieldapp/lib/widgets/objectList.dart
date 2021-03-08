import 'package:fieldapp/main.dart';
import 'package:fieldapp/types/types.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

class ObjectList extends StatefulWidget {
  final String header;
  final List<GenesisObject> list;
  final Function(int) dismissItem;
  final String highlightUUID;
  ObjectList({
    Key key,
    @required this.header,
    @required this.list,
    @required this.dismissItem,
    @required this.highlightUUID,
  }) : super(key: key);

  @override
  _ObjectListState createState() => _ObjectListState();
}

class _ObjectListState extends State<ObjectList> {
  @override
  Widget build(BuildContext context) {
    return Container(
      child: Column(
        children: <Widget>[
          // Header
          AnimatedContainer(
            duration: Duration(milliseconds: 250),
            height: widget.list.length == 0 ? 0 : 28,
            width: double.infinity,
            decoration: BoxDecoration(color: COLOUR_PRIMARY),
            child: Center(
              child: Text(
                widget.header + (widget.list.length > 1 ? "s" : ""),
                textAlign: TextAlign.center,
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 18,
                ),
              ),
            ),
          ),
          // List
          ListView.builder(
            shrinkWrap: true,
            physics: NeverScrollableScrollPhysics(),
            itemBuilder: (context, i) => Dismissible(
              key: Key(widget.list[i].id),
              background: Container(
                decoration: BoxDecoration(
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black,
                      offset: const Offset(0.0, 0.0),
                    ),
                    const BoxShadow(
                      color: COLOUR_PRIMARY,
                      blurRadius: 6,
                      spreadRadius: 0,
                      offset: const Offset(0, 0),
                    )
                  ],
                ),
              ),
              onDismissed: (dir) => widget.dismissItem(i),
              child: codeItem(
                i,
                widget.list[i],
                color: widget.highlightUUID == widget.list[i].id
                    ? Colors.red.shade400
                    : i % 2 == 0 ? Colors.white : Colors.grey.shade100,
              ),
            ),
            itemCount: widget.list.length,
          ),
        ],
      ),
    );
  }

  Widget codeItem(
    int index,
    GenesisObject obj, {
    Color color,
  }) {
    List<Widget> children = [];
    if (obj is Product) {
      if (obj.order != null) children.add(codeItemSmall(obj.order));
      if (obj.carton != null) children.add(codeItemSmall(obj.carton));
      if (obj.sku != null) children.add(codeItemSmall(obj.sku));
    } else if (obj is Carton) {
      if (obj.order != null) children.add(codeItemSmall(obj.order));
      if (obj.pallet != null) children.add(codeItemSmall(obj.pallet));
    } else if (obj is Pallet && obj.container != null) {
      children.add(codeItemSmall(obj.container));
    }

    BoxDecoration rowBorder() {
      return BoxDecoration(
        border: Border(bottom: BorderSide(color: Colors.black, width: 2)),
      );
    }

    return AnimatedContainer(
      duration: Duration(milliseconds: 200),
      color: color,
      child: Container(
        decoration: rowBorder(),
        child: Padding(
          child: Row(
            children: <Widget>[
              Container(
                margin: EdgeInsets.only(top: 5),
                child: SizedBox(
                  child: Text(
                    (index + 1).toString() + ".",
                    style: TextStyle(fontSize: 20),
                  ),
                  width: 30,
                ),
              ),
              SizedBox(
                child: FaIcon(obj.icon),
                width: 40,
              ),
              Expanded(
                child: Column(
                  children: <Widget>[
                    Text(
                      obj.code.length > 14
                          ? obj.code.substring(0, 13) + "..."
                          : obj.code,
                      style: TextStyle(fontSize: 16),
                    ),
                    Row(
                      children: children,
                    ),
                  ],
                  crossAxisAlignment: CrossAxisAlignment.start,
                ),
              ),
              SizedBox(
                child: IconButton(
                  icon: Icon(Icons.delete),
                  tooltip: "delete item",
                  onPressed: () {
                    setState(() {
                      widget.dismissItem(index);
                    });
                  },
                ),
                width: 40,
              ),
            ],
          ),
          padding: EdgeInsets.symmetric(vertical: 8, horizontal: 10),
        ),
      ),
    );
  }

  Widget codeItemSmall(GenesisObject obj, {Color color, IconData icon}) {
    return Padding(
      child: Row(
        children: <Widget>[
          Padding(
            child: FaIcon(
              obj.icon,
              size: 12,
            ),
            padding: EdgeInsets.only(right: 4),
          ),
          Text(
            obj.code,
            style: TextStyle(fontSize: 10),
          ),
        ],
      ),
      padding: EdgeInsets.only(right: 8),
    );
  }
}
