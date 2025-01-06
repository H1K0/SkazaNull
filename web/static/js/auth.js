$("#auth").on("submit", function (e) {
	e.preventDefault();
	$.ajax({
		url: "/api/auth",
		type: "POST",
		data: $("#auth").serialize(),
		dataType: "json",
		success: function (resp) {
			location.reload();
		},
		error: function (err) {
			$("#error-message").text(err.responseJSON.error);
			$("#error").removeClass("hidden");
		},
	});
});
