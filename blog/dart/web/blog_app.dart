import "dart:html";
import "dart:convert" show JSON;
import 'dart:async';
import "package:polymer/polymer.dart";

@CustomTag('blog-app')
class blog_app extends PolymerElement {
  @observable List<Map> posts;

  /// For page scroll position
  DivElement page;
  String prefix = "kunblog";

  /// For HTML history
  Map stateAll = {"post": "all"};
  Map stateContent = {"post": "content"};

  blog_app.created() : super.created();

  void attached() {
    super.attached();
    page = this.shadowRoot.querySelector("#pages");
    loadAllPostMeta();
  }

  detached() {
  }

  void loadAllPostMeta() {
    //String url = window.location.href + "all";
    String url = "http://likunarmstrong.appspot.com/blog/all";
    HttpRequest.getString(url).then(onPostMetaDataLoaded);
  }

  void resumeWindowPosition() {
    String name = window.name;
    print("Window name is $name");
    if (name.contains(prefix)) {
      int position = int.parse(name.split("_")[1], onError:(String source) => 0);
      print("resume to position $position");
      page.scrollTop = position;
    }
  }

  void onPostMetaDataLoaded(String responseText) {
    if (responseText.length == 0) {
      return;
    }
    posts = JSON.decode(responseText);
    print("Returned json is : $posts");

    new Timer(new Duration(milliseconds:500), () => resumeWindowPosition());
  }

  /// Callback
  void selectPost(Event e, var detail, Node target) {
    e.preventDefault();
    window.name = prefix + "_${$["pages"].scrollTop}";

    // Prepare data
    String id = (target as DivElement).id;
    int index = int.parse(id.split("-")[1]);
    String link = posts[index]["link"];
    String fileName = posts[index]["filename"];
    if (link.contains("http")) {  // Redirect if it's external pdf
      window.location.assign(link);
    } else {  // Push history state
      window.location.assign(window.location.href + "/" + fileName);
    }
  }
}