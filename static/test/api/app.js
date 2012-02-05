hiro.module('APP', {
  onTest: function() {
    var app = { name: 'coma' };
    var html = '<html><head><title>COMA</title></head><body></body></html>';

    this.args = [ app, html ];
  },
  "test PUT": function(app) {
    var h = this;

    h.expect(1);
    h.pause();

    PUT({
      url: '/' + app.name,
      data: app
    }, function(res) {
      h.assertEqual(res.status, 204);

      h.resume();
    });
  },
  "test PUT static html": function(app, html) {
    var h = this;

    h.expect(1);
    h.pause();

    PUT({
      url: '/' + app.name + '/index.html',
      data: html
    }, function(res) {
      h.assertEqual(res.status, 204);

      h.resume();
    });
  },
  "test GET static html": function(app, html) {
    var h = this;

    h.expect(2);
    h.pause();

    GET({
      url: '/' + app.name + '/index.html',
    }, function(res) {
      h.assertEqual(res.status, 200);
      h.assertEqual(res.body, html)

      h.resume();
    });
  },
});
