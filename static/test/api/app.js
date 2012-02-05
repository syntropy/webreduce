hiro.module('APP', {
  onTest: function() {
    var app = { name: 'coma' };

    this.args = [ app ];
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
  "test PUT static html": function(app) {
    var h = this;

    h.expect(1);
    h.pause();

    PUT({
      url: '/' + app.name + '/index.html',
      data: '<html><head><title>COMA</title></head><body></body></html>'
    }, function(res) {
      h.assertEqual(res.status, 204);

      h.resume();
    });
  },
});
