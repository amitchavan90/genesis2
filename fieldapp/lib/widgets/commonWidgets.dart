import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/widgets.dart';

Widget basicTextField(
  String name,
  String defaultValue,
  Map<String, TextEditingController> controllers,
  Map<String, FocusNode> focusNodes, {
  TextInputType keyboardType = TextInputType.text,
  int maxLines = 1,
  bool obscureText = false,
  Function(String) validator,
  TextInputAction textInputAction,
  Function(String) onFieldSubmitted,
  Function(String) onChanged,
}) {
  if (!controllers.containsKey(name)) {
    controllers[name] = TextEditingController(text: defaultValue);
  }
  if (!focusNodes.containsKey(name)) focusNodes[name] = FocusNode();

  return Padding(
    padding: EdgeInsets.only(bottom: 15),
    child: TextFormField(
      controller: controllers[name],
      focusNode: focusNodes[name],
      validator: validator,
      keyboardType: keyboardType,
      obscureText: obscureText,
      textInputAction: textInputAction,
      maxLines: maxLines,
      decoration: InputDecoration(
        labelText: name,
        filled: true,
        fillColor: Colors.white.withOpacity(0.1),
        labelStyle: TextStyle(
          color: Color(0xFFA1A1A1),
        ),
        errorStyle: TextStyle(
          height: 0,
          backgroundColor: Colors.red.shade800.withOpacity(0.9),
          color: Colors.white,
          fontSize: 9,
        ),
        enabledBorder: UnderlineInputBorder(
          borderSide: BorderSide(width: 1.5),
        ),
        focusedBorder: UnderlineInputBorder(
          borderSide: BorderSide(width: 1.5, color: Color(0xFF009DFF)),
        ),
        errorBorder: UnderlineInputBorder(
          borderSide: BorderSide(width: 2, color: Colors.red.shade800),
        ),
        contentPadding: EdgeInsets.only(left: 5, right: 5, bottom: 5),
      ),
      onFieldSubmitted: (value) {
        if (textInputAction == TextInputAction.next) {
          int i = focusNodes.entries.toList().indexWhere((e) => e.key == name);
          if (focusNodes.entries.length > i + 1) {
            focusNodes.entries.elementAt(i + 1).value.requestFocus();
          }
        }
        if (onFieldSubmitted != null) onFieldSubmitted(value);
      },
      onChanged: onChanged,
    ),
  );
}

Widget textFieldWithIcon(
  String name, {
  TextInputType keyboardType,
  TextInputAction textInputAction,
  bool obscureText = false,
  TextEditingController controller,
  String Function(String) validator,
  Function(String) onFieldSubmitted,
  Function(String) onChanged,
  FocusNode focusNode,
  FocusNode focusNext,
  BuildContext context,
  Widget icon,
  Widget iconExtra,
  double contentPaddingRight = 0,
}) {
  return Stack(
    children: <Widget>[
      icon != null
          ? Padding(
              padding: EdgeInsets.only(left: 2, top: 10),
              child: icon,
            )
          : Container(),
      TextFormField(
        keyboardType: keyboardType,
        controller: controller,
        textInputAction: textInputAction,
        focusNode: focusNode,
        obscureText: obscureText,
        decoration: InputDecoration(
          labelText: name,
          labelStyle: TextStyle(color: Color(0xFFA1A1A1)),
          errorStyle: TextStyle(
            backgroundColor: Colors.red.shade800.withOpacity(0.9),
            color: Colors.white,
            fontSize: 12,
          ),
          enabledBorder: UnderlineInputBorder(
            borderSide: BorderSide(width: 1.5),
          ),
          focusedBorder: UnderlineInputBorder(
            borderSide: BorderSide(width: 1.5, color: Color(0xFF009DFF)),
          ),
          errorBorder: UnderlineInputBorder(
            borderSide: BorderSide(width: 2, color: Colors.red.shade800),
          ),
          contentPadding: EdgeInsets.only(
            bottom: 4,
            left: icon != null ? 32 : 0,
            right: contentPaddingRight,
          ),
        ),
        validator: validator,
        onFieldSubmitted: (value) {
          if (focusNext != null) FocusScope.of(context).requestFocus(focusNext);
          if (onFieldSubmitted != null) onFieldSubmitted(value);
        },
        onChanged: onChanged,
      ),
      iconExtra != null
          ? Positioned(
              right: 2,
              top: 10,
              child: iconExtra,
            )
          : Container()
    ],
  );
}

Widget basicTextButton({
  String text,
  Widget child,
  double height = 40,
  double width,
  Color color,
  Color color2,
  Color textColor = Colors.white,
  Function onTap,
}) {
  return Container(
    height: height,
    width: width,
    child: Stack(
      children: <Widget>[
        Container(
          color: color2 != null ? null : color,
          decoration: color2 != null
              ? BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.centerLeft,
                    end: Alignment.centerRight,
                    stops: [0, 1],
                    colors: [color, color2],
                  ),
                )
              : null,
        ),
        Material(
          type: MaterialType.transparency,
          child: InkWell(
            onTap: onTap,
            child: child ??
                Center(
                  child: Text(
                    text,
                    style: TextStyle(
                      color: textColor,
                      fontWeight: FontWeight.w600,
                      fontSize: 16,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ),
            splashColor: Colors.white.withOpacity(0.4),
            highlightColor: Colors.white.withOpacity(0.2),
          ),
        ),
      ],
    ),
  );
}

Widget loadingOverlay(bool visible) {
  return Positioned.fill(
    child: IgnorePointer(
      ignoring: !visible,
      child: Stack(
        children: [
          BackdropFilter(
            filter: ui.ImageFilter.blur(sigmaX: 2, sigmaY: 2),
            child: AnimatedOpacity(
              opacity: visible ? 0.4 : 0,
              duration: Duration(milliseconds: 400),
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.black,
                ),
              ),
            ),
          ),
          AnimatedOpacity(
            opacity: visible ? 1 : 0,
            duration: Duration(milliseconds: 400),
            child: Center(
              child: CircularProgressIndicator(strokeWidth: 6),
            ),
          ),
        ],
      ),
    ),
  );
}
