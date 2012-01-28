hiro.module('Agent', {
  onTest: function() {
    var agent = {
      code: 'local params = {...}; emit(params[0]);',
      language: 'lua',
      name: 'foo',
      url: '/agents/foo'
    };

    this.args = [ agent ];
  },
  "test GET agents": function() {
    var h = this;

    h.expect(2);
    h.pause();


    GET({ url: '/agents' }, function(res) {
      h.assertEqual(res.status, 200);
      h.assertEqual(res.body.result.length, res.body.count);

      h.resume();
    });
  },
  "test PUT agent": function(agent) {
    var h = this;

    h.expect(1);
    h.pause();

    PUT({ url: agent.url, data: { code: agent.code, language: agent.language } }, function(res) {
      h.assertEqual(res.status, 204);

      h.resume()
    });
  },
  "test GET agent": function(agent) {
    var h = this;

    h.expect(4);
    h.pause();

    GET({ url: agent.url }, function(res) {
      h.assertEqual(res.status, 200);
      h.assertEqual(res.body.name, agent.name);
      h.assertEqual(res.body.code, agent.code);
      h.assertEqual(res.body.language, agent.language);

      h.resume();
    });
  },
  "test POST agent": function(agent) {
    var h = this;

    h.expect(1);
    h.pause();

    POST({ url: agent.url, data: { foo: 'bar' } }, function(res) {
      h.assertEqual(res.status, 202);

      h.resume();
    });
  },
  "test POST on missing agent": function() {
    var h = this;

    h.expect(1);
    h.pause();

    POST({ url: '/agents/bar', data: { foo: 'bar' } }, function(res) {
      h.assertEqual(res.status, 404);

      h.resume();
    });
  }
});
