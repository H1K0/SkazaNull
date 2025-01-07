$(window).on("load", function (e) {
	$.ajax({
		url: "/api/auth",
		type: "GET",
		dataType: "json",
		success: function (resp) {
			$("#input-name").val(resp.name);
			$("#input-login").val(resp.login);
			$("#input-tgid").val(resp.telegram_id);
		},
		error: function (err) {
			$("#error-message").text(err.responseJSON.error);
			$("#error").removeClass("hidden");
		},
	});
});

$(document).on("click", "#btn-logout", function (e) {
	$.ajax({
		url: "/api/auth",
		type: "DELETE",
		success: function () {
			location.reload();
		},
		error: function (err) {
			$("#error-message").text(err.responseJSON.error);
			$("#error").removeClass("hidden");
		},
	});
});

$(document).on("submit", "#user-update", function (e) {
	e.preventDefault();
	data = formToJSON($("#user-update"));
	$.ajax({
		url: "/api/auth",
		type: "PATCH",
		contentType: "application/json",
		data: JSON.stringify(data),
		processData: false,
		dataType: "json",
		success: function () {
			$("#error").addClass("hidden");
			$("#success").removeClass("hidden");
            $("#input-password").val("");
		},
		error: function (err) {
			$("#success").addClass("hidden");
			$("#error-message").text(err.responseJSON.error);
			$("#error").removeClass("hidden");
		},
	});
});
