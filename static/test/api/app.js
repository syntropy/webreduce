hiro.module('APP', {
  onTest: function() {
    var app = {
      name: 'coma',
      html: '<html><head><title>COMA</title></head><body></body></html>'
    };
    var sensor = {
      name: 'random'
    };

    this.args = [ app, sensor ];
  },
  "test PUT": function(app) {
    var h = this;

    h.expect(1);
    h.pause();

    PUT({
      url: '/' + app.name,
      data: { name: app.name }
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
      data: app.html
    }, function(res) {
      h.assertEqual(res.status, 204);

      h.resume();
    });
  },
  "test GET static html": function(app) {
    var h = this;

    h.expect(2);
    h.pause();

    GET({
      url: '/' + app.name + '/index.html',
    }, function(res) {
      h.assertEqual(res.status, 200);
      h.assertEqual(res.body, app.html)

      h.resume();
    });
  },
  "test PUT app sensor": function(app, sensor) {
    var h = this;

    h.expect(1);
    h.pause();

    PUT({
      url: '/' + app.name + '/sensors/' + sensor.name,
      data: { name: sensor.name, sensor: "POST" }
    }, function(res) {
      h.assertEqual(res.status, 204);

      h.resume();
    });
  },
  "test POST app signal": function(app) {
    var h = this;

    h.expect(1);
    h.pause();

    POST({
      url: '/' + app.name,
      data: { action: 'start' }
    }, function(res) {
      h.assertEqual(res.status, 202);

      h.resume();
    });
  },
  "test POST to WebSocket": function(app, sensor) {
    var h = this;
    var path = 'ws://' + location.host + '/' + app.name + '/ws';
    var ws = new WebSocket(path);

    h.expect(2);
    h.pause();

    ws.onmessage = function() {
      ws.close();

      h.assertTrue(true);
      h.resume();
    }

    setTimeout(function() {
      POST({
        url: '/' + app.name + '/sensors/' + sensor.name,
        data: { foo: 'bla' }
      }, function(res) {
        h.assertEqual(res.status, 202);
      });
    }, 1)
  },
  "test DELETE app sensor": function(app, sensor) {
    var h = this;

    h.expect(2);
    h.pause();

    DELETE({
      url: '/' + app.name + '/sensors/' + sensor.name,
      data: { name: sensor.name, sensor: "POST" }
    }, function(res) {
      h.assertEqual(res.status, 204);

      h.resume();
    });

    GET({
      url: '/' + app.name + '/sensors/' + sensor.name
    }, function(res) {
      h.assertEqual(res.status, 404);

      h.resume();
    });
  }
});
