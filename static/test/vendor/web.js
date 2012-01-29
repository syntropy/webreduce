/*global hiro:false, jQuery: false, $:false */

(function ($, undefined) {
	"use strict";

	var document = window.document;
	var context  = '#web';

	hiro.bind('hiro.onComplete', function () {
		$(document.createElement('p'))
			.addClass('simple')
			.html('All tests finished')
			.appendTo(context);
	});

	hiro.bind('suite.onStart', function (suite) {
		var uid = 'hiro_suite_' + suite.name;
		var div = document.createElement('div');

		$(div)
			.addClass('suite')
			.addClass('idle')
			.attr('id', uid);

		$(document.createElement('h2'))
			.html(suite.name)
			.appendTo(div);

		$(context).append(div);

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
		var div = document.createElement('div');

		$(div)
			.addClass('test')
			.addClass('idle')
			.attr('id', uid)
			.html(test.name);

		$(document.createElement('div'))
			.addClass('report')
			.hide()
			.appendTo(div);

		$(div).appendTo(context);
		context = '#' + uid;
	});

	hiro.bind('test.onFailure', function (test, report) {
		var div = $('div.report', context);

		$(document.createElement('p'))
			.html('<label>Assertion:</label> ' + report.assertion)
			.appendTo(div);

		if (report.expected) {
			$(document.createElement('p'))
				.html('<label>Expected:</label> ' + report.expected)
				.appendTo(div);
		}

		$(document.createElement('p'))
			.html('<label>Result:</label> ' + report.result)
			.appendTo(div);

		$(document.createElement('p'))
			.html('<label>Position:</label> ' + report.position)
			.appendTo(div);
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
