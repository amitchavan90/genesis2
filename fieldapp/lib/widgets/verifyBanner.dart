import 'package:fieldapp/auth.dart';
import 'package:fieldapp/main.dart';
import 'package:fieldapp/types/types.dart';
import 'package:flutter/material.dart';
import 'package:fieldapp/widgets/dialogs.dart';

class VerifyBanner extends StatefulWidget {
  final String message;
  final Color colour;
  final void Function(bool) setLoading;
  VerifyBanner({
    Key key,
    this.message,
    this.colour,
    @required this.setLoading,
  }) : super(key: key);

  @override
  _VerifyBannerState createState() => _VerifyBannerState();
}

class _VerifyBannerState extends State<VerifyBanner> {
  final formKey = GlobalKey<FormState>();

  void onTap() async {
    widget.setLoading(true);
    GQLResponse response = await Auth.forgotPassword(me.email);
    widget.setLoading(false);

    if (!response.success) {
      showErrorDialog(context, response.message);
      return;
    }

    showSuccessDialog(
      context,
      'Code sent',
      'A verification code has been sent to ${me.mobilePhone}',
    ).then(
      (value) => showEnterCodeDialog(
        context,
        formKey,
        (String token) async {
          GQLResponse verifyResponse =
              await Auth.verifyResetToken(token, me.email);

          if (!verifyResponse.success) {
            showErrorDialog(context, verifyResponse.message);
            return false;
          }

          // Change password dialog
          Navigator.of(context).pop();
          showSuccessDialog(context, "Verification successful!", "");
          me.mobileVerified = true;

          return false; // don't pop context (already popped)
        },
      ),
    );
  }

  Widget build(BuildContext context) {
    return Stack(
      children: <Widget>[
        Container(
          height: 30,
          decoration: BoxDecoration(
            color: widget.colour != null ? widget.colour : Colors.green,
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.3),
                blurRadius: 3,
                offset: Offset(0, 2),
              )
            ],
          ),
        ),
        Material(
          type: MaterialType.transparency,
          child: InkWell(
            child: Container(
              height: 30,
              child: Center(
                child: Text(
                  widget.message,
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ),
            onTap: onTap,
          ),
        ),
      ],
    );
  }
}
