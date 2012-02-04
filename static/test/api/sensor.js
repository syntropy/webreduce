hiro.module('Sensor', {
  onTest: function() {
    var sensor = {
      code: 'local params = {...}; emit("\\\""..params[1].."\\\" ← this is emmited data!");',
      language: 'lua',
      name: 'foo',
      url: '/sensors/foo'
    };

    this.args = [ sensor ];
  },
  "test GET sensors": function() {
    var h = this;

    h.expect(2);
    h.pause();


    GET({ url: '/sensors' }, function(res) {
      h.assertEqual(res.status, 200);
      h.assertEqual(res.body.items.length, res.body.count);

      h.resume();
    });
  },
  "test PUT sensor": function(sensor) {
    var h = this;

    h.expect(1);
    h.pause();

    PUT({ url: sensor.url, data: { code: sensor.code, language: sensor.language } }, function(res) {
      h.assertEqual(res.status, 204);

      h.resume()
    });
  },
  "test GET sensor": function(sensor) {
    var h = this;

    h.expect(4);
    h.pause();

    GET({ url: sensor.url }, function(res) {
      h.assertEqual(res.status, 200);
      h.assertEqual(res.body.name, sensor.name);
      h.assertEqual(res.body.code, sensor.code);
      h.assertEqual(res.body.language, sensor.language);

      h.resume();
    });
  },
  "test POST sensor": function(sensor) {
    var h = this;

    h.expect(1);
    h.pause();

    POST({ url: sensor.url, data: { foo: 'bar' } }, function(res) {
      h.assertEqual(res.status, 202);

      h.resume();
    });
  },
  "test PUT with invalid code": function(sensor) {
    var h = this;

    h.expect(1);
    h.pause();

    PUT({ url: sensor.url, data: { code: 'that is no code, no good', language: sensor.language } }, function(res) {
      h.assertEqual(res.status, 400);

      h.resume()
    });
  },
  "test POST on missing sensor": function() {
    var h = this;

    h.expect(1);
    h.pause();

    POST({ url: '/sensors/bar', data: { foo: 'bar' } }, function(res) {
      h.assertEqual(res.status, 404);

      h.resume();
    });
  }
});
