hiro.module('Agent', {
  setUp: function() {},
  "test GET agents": function() {
    var h = this;

    h.expect(2);
    h.pause();

    jQuery.ajax({
      url: '/agents',
      type: 'get',
      dataType: 'json',
      complete: function(res) {
        var payload = JSON.parse(res.responseText);
        h.resume();

        h.assertEqual(res.status, 200);
        h.assertEqual(payload.result.length, payload.count);
      }
    });
  },
  "test PUT and GET agent": function() {
    var h = this;
    var data = { language: 'lua', code: 'local params = {...}; emit(params[0]);' };
    var name = 'foo';
    var url = '/agents/' + name;

    h.expect(5);
    h.pause();

    jQuery.ajax({
      url: url,
      type: 'put',
      dataType: 'json',
      contentType: 'application/json',
      data: JSON.stringify(data),
      processData: false,
      complete: function(res) {
        h.assertEqual(res.status, 204);

        jQuery.ajax({
          url: url,
          type: 'get',
          dataType: 'json',
          complete: function(res) {
            var payload = JSON.parse(res.responseText);

            h.assertEqual(res.status, 200);
            h.assertEqual(payload.name, name);
            h.assertEqual(payload.code, data.code);
            h.assertEqual(payload.language, data.language);

            h.resume();
          }
        });
      }
    });
  }
});
