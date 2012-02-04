/*global hiro:false, jQuery: false, $:false */

(function ($, undefined) {
	"use strict";

	var context  = '#web';

	hiro.bind('hiro.onComplete', function () {
		$('<p>', {
		  'class': 'simple',
		  'text': 'All tests finished'
    }).appendTo(context);
	});

	hiro.bind('suite.onStart', function (suite) {
		var uid = 'hiro_suite_' + suite.name;

    $('<div>', {
      'class': 'idle suite',
      'id': uid
    }).appendTo(context)
      .append($('<h2>', { 'text': suite.name }));

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
		var uid = 'hiro_test_' + test.suite.name + '_' + test.name;

    $('<div>', {
      'class': 'idle test',
      'id': uid,
      'html': test.name
    }).appendTo(context)
      .append($('<div>', { 'class': 'report' }).hide());

		context = '#' + uid;
	});

	hiro.bind('test.onFailure', function (test, report) {
		var div = $('div.test').filter(function(i, elem) {
		  return $(this).text() === test.name;
		}).children('div.report');

    div.append($('<p>', {
      html: '<label>Assertion:</label> ' + report.assertion
    }));

		if (report.expected) {
      div.append($('<p>', {
        html: '<label>Expected:</label> ' + report.expected
      }));
		}

    div.append($('<p>', {
      html: '<label>Result:</label> ' + report.result
    }));

    div.append($('<p>', {
      html: '<label>Position:</label> ' + report.position
    }));
	});

	hiro.bind('test.onComplete', function (test, success) {
	  console.dir(arguments)
		$(context)
			.removeClass('idle')
			.addClass(success ? 'succ' : 'fail');

		if (!success)
		  console.log(context)
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
