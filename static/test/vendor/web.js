/*
 * Copyright (c) 2011 Anton Kovalyov, http://hirojs.com/
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the
 * "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use, copy, modify, merge, publish,
 * distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to
 * the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 * LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
 * OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

/*global hiro:false, ender:false */

(function ($, window, undefined) {
  "use strict";

  var context  = '#web';
  var start;

  hiro.bind('hiro.onStart', function() {
    start = Date.now();
  });

  hiro.bind('hiro.onComplete', function () {
    var duration = Date.now() - start;

    $('<div>', {
      'class': 'result succ',
      'html': '(' + duration + 'ms)'
    }).appendTo(context);
  });

  hiro.bind('suite.onStart', function (suite) {
    var uid = 'hiro_suite_' + suite.name;

    $('<div>', {
      'class': 'idle suite',
      'id': uid
    }).append($('<h2>', {
      'html': '<a href="?' + suite.name + '">' + suite.name + '</a>'
    })).appendTo(context);

    context = '#' + uid;
  });

  hiro.bind('suite.onComplete', function (suite, success) {
    $(context)
      .removeClass('idle')
      .addClass(success ? 'succ' : 'fail');

    context = '#web';
  });

  hiro.bind('suite.onTimeout', function (suite, success) {
    $(context)
      .removeClass('idle')
      .addClass('fail');

    context = '#web';
  });

  hiro.bind('test.onStart', function (test) {
    var name = test.name.replace(/^test/, '');
    var uid = nextId();

    $('<div>', {
      'class': 'idle test',
      'id': uid,
      'html': '<span>' + name + '</span>'
    }).append($('<div>', {
      'class': 'report'
    })).appendTo(context);

    context = '#' + uid;
  });

  hiro.bind('test.onFailure', function (test, report) {
    var div = $('div.report', context);

    $('<p>', {
      'html': '<label>Assertion:</label> ' + report.assertion
    }).appendTo(div)

    if (report.expected) {
      $('<p>', {
        'html': '<label>Expected:</label> ' + report.expected
      }).appendTo(div);
    }

    $('<p>', {
      'html': '<label>Result:</label> ' + report.result
    }).appendTo(div);

    $('<p>', {
      'html': '<label>Position:</label> ' + report.position
    }).appendTo(div);
  });

  hiro.bind('test.onComplete', function (test, success) {
    $(context)
      .removeClass('idle')
      .addClass(success ? 'succ' : 'fail');

    if (!success)
      $('div.report', context).show();

    context = '#hiro_suite_' + test.suite.name;
  });

  hiro.bind('test.onTimeout', function (test, success) {
    $(context)
      .removeClass('idle')
      .addClass('fail');

    context = '#hiro_suite_' + test.suite.name;
  });
}(jQuery, window));

var id = 0;

function nextId() {
  return id++;
}
