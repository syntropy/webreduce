function req(cb, opts) {
  var opts = jQuery.extend({}, {
    headers: {
      accept: 'application/json'
    },
    accepts: 'json',
    dataType: 'json',
    contentType: 'application/json',
    processData: false,
    complete: function(res) {
      if (res.responseText.length) {
        res.body = JSON.parse(res.responseText);
      } else {
        res.body = res.responseText;
      }

      cb(res);
    }
  }, opts);

  jQuery.ajax(opts);
}

function GET(opts, cb) {
  req(cb, jQuery.extend({}, { type: 'get' }, opts));
}

function PUT(opts, cb) {
  req(cb, jQuery.extend({}, { type: 'put' }, opts));
}

function POST(opts, cb) {
  req(cb, jQuery.extend({}, { type: 'post' }, opts));
}
function DELETE(opts, cb) {
  req(cb, jQuery.extend({}, { type: 'delete' }, opts));
}
