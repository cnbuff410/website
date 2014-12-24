import 'dart:html';
import 'package:polymer/polymer.dart';

/**
 * A Polymer click counter element.
 */
@CustomTag('post-card')
class PostCard extends PolymerElement with Polymer, Observable {
  @published String title;
  @published String date;
  @observable int shadowLevel = 1;
  @observable String year = "1984";
  @observable String month = "五月";
  @observable String day = "24";
  @observable bool smallScreen;

  PostCard.created() : super.created();
  @override
  void attached() {
    super.attached();
    List<String> dateFields = date.split("-");
    if (dateFields.length == 3) {
      year = dateFields[0];
      month = dateFields[1] + "月";
      day = dateFields[2] + "日";
    }
  }

  @override
  void detached() {
    super.detached();
  }

  void onTabClicked(Event e, var detail, Node target) {
  }

  void mouseOver(Event e, var detail, Node target) {
    shadowLevel = 3;
  }

  void mouseLeave(Event e, var detail, Node target) {
    shadowLevel = 1;
  }
}
