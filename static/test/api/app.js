hiro.module('APP', {
  onTest: function() {
    var app = { name: 'caro' };

    this.args = [ app ];
  },
  "test PUT": function(app) {
    var h = this;

    h.expect(1);
    h.pause();

    PUT({ url: '/' + app.name, data: app }, function(res) {
      h.assertEqual(res.status, 200);

      h.resume();
    });
  },
});
