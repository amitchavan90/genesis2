import 'package:fieldapp/types/enums.dart';
import 'package:fieldapp/utils.dart';
import 'package:flutter/widgets.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

class GQLResponse {
  bool success = false;
  String message = "An issue occurred.";
}

class User {
  String id;
  String email;
  String firstName;
  String lastName;
  String mobilePhone;
  Role role;

  bool verified;
  bool mobileVerified;

  User.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        email = json['email'],
        firstName = json['firstName'],
        lastName = json['lastName'],
        mobilePhone = json['mobilePhone'],
        role = Role.fromJson(json['role']),
        verified = json['verified'],
        mobileVerified = json['mobileVerified'];
}

class Role {
  String id;
  String name;
  int tier;
  List<Perm> permissions = [];
  List<TrackAction> trackActions = [];

  Role.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        name = json['name'],
        tier = json['tier'],
        permissions = StringToEnum.fromList(
          Perm.values,
          (json['permissions'] as List).cast<String>().toList(),
        ),
        trackActions = (json['trackActions'] as List)
            .map((e) => TrackAction.fromJson(e))
            .toList();
}

class GenesisObject {
  String id;
  String code;
  bool archived;

  IconData icon;

  GenesisObject.fromJson(Map<String, dynamic> json, {IconData icon})
      : id = json['id'],
        code = json['code'] == null ? json['id'] : json['code'],
        archived = json['archived'],
        icon = icon ?? FontAwesomeIcons.barcodeAlt;
}

class SKU extends GenesisObject {
  String name;
  String description;
  String masterPlanURL;

  SKU.fromJson(Map<String, dynamic> json)
      : name = json['name'],
        description = json['description'],
        masterPlanURL =
            json['masterPlan'] != null ? json['masterPlan']['file_url'] : null,
        super.fromJson(json);
}

class Order extends GenesisObject {
  Order.fromJson(Map<String, dynamic> json)
      : super.fromJson(json, icon: FontAwesomeIcons.shoppingCart);
}

class Product extends GenesisObject {
  int loyaltyPoints;
  Carton carton;
  SKU sku;
  Order order;

  Product.fromJson(Map<String, dynamic> json)
      : loyaltyPoints = json['loyaltyPoints'],
        carton =
            json['carton'] != null ? Carton.fromJson(json['carton']) : null,
        sku = json['sku'] != null ? SKU.fromJson(json['sku']) : null,
        order = json['order'] != null ? Order.fromJson(json['order']) : null,
        super.fromJson(json, icon: FontAwesomeIcons.lightSteak);
}

class Carton extends GenesisObject {
  int productCount;
  Pallet pallet;
  Order order;
  Carton.fromJson(Map<String, dynamic> json)
      : productCount =
            json['productCount'] != null ? json['productCount'] : null,
        pallet =
            json['pallet'] != null ? Pallet.fromJson(json['pallet']) : null,
        order = json['order'] != null ? Order.fromJson(json['order']) : null,
        super.fromJson(json, icon: FontAwesomeIcons.box);
}

class Pallet extends GenesisObject {
  int cartonCount;
  GenesisContainer container;

  Pallet.fromJson(Map<String, dynamic> json)
      : cartonCount = json['cartonCount'],
        container = json['container'] != null
            ? GenesisContainer.fromJson(json['container'])
            : null,
        super.fromJson(json, icon: FontAwesomeIcons.palletAlt);
}

class GenesisContainer extends GenesisObject {
  int palletCount;

  GenesisContainer.fromJson(Map<String, dynamic> json)
      : palletCount = json['palletCount'],
        super.fromJson(json, icon: FontAwesomeIcons.containerStorage);
}

class TrackAction extends GenesisObject {
  String name;
  String nameChinese;
  List<bool> requirePhotos;

  TrackAction.fromJson(Map<String, dynamic> json)
      : name = json['name'],
        nameChinese = json['nameChinese'],
        requirePhotos = List.castFrom(json['requirePhotos']),
        super.fromJson(json, icon: FontAwesomeIcons.solidTruckMoving);
}
