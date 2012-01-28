hiro.module('Agent', {
  setUp: function() {},
  "test GET agents": function() {
    var h = this;

    h.expect(2);
    h.pause();

    GET({ url: '/agents' }, function(res) {
      var payload = JSON.parse(res.responseText);
      h.resume();

      h.assertEqual(res.status, 200);
      h.assertEqual(payload.result.length, payload.count);
    });
  },
  "test PUT and GET agent": function() {
    var h = this;
    var data = { language: 'lua', code: 'local params = {...}; emit(params[0]);' };
    var name = 'foo';
    var url = '/agents/' + name;

    h.expect(5);
    h.pause();

    PUT({ url: url, data: JSON.stringify(data) }, function(res) {
      h.assertEqual(res.status, 204);

      GET({ url: url }, function(res) {
        h.assertEqual(res.status, 200);
        h.assertEqual(res.body.name, name);
        h.assertEqual(res.body.code, data.code);
        h.assertEqual(res.body.language, data.language);

        h.resume();
      });
    });
  }
});
