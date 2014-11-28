import "dart:html";
import "dart:convert" show JSON;
import "package:polymer/polymer.dart";

@CustomTag('post-app')
class post_app extends PolymerElement {
  post_app.created() : super.created();

  void attached() {
    super.attached();
    loadContent(window.location.href + "?mode=raw");
  }

  detached() {
  }

  void loadContent(String postUrl) {
    HttpRequest.getString(postUrl).then(onPostContentDataLoaded);
  }

  void onPostContentDataLoaded(String responseText) {
    if (responseText.length == 0) {
      return;
    }
    print("response is ${responseText}");
    DivElement injectDiv = this.shadowRoot.querySelector('#inject');
    mInjectBoundHTML(responseText, injectDiv);
  }

  final NodeValidatorBuilder _htmlValidator = new NodeValidatorBuilder.common()
      ..allowElement('a', attributes: ['href'])
      ..allowElement('img', attributes: ['src']);

  DocumentFragment mInjectBoundHTML(String html, [Element element]) {
    var template =
        new TemplateElement()..setInnerHtml(html, validator: _htmlValidator);
    var fragment = this.instanceTemplate(template);
    if (element != null) {
      element.text = '';
      element.append(fragment);
    }
    return fragment;
  }

  /// Callback
  void back(Event e, var detail, Node target) {
    e.preventDefault();
    window.history.back();
  }
}
