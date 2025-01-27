const PAGE_SIZE = 10;
var totalPages;
var currPage = +sessionStorage.getItem("page");
if (currPage == 0) {
	currPage = 1;
}
var search = sessionStorage.getItem("search");
if (search == null) {
	search = "";
}
var sorting = sessionStorage.getItem("sort");
if (sorting == null) {
	sorting = "-datetime";
}
var temp_quote_id;

function renderBlockQuote(quote) {
	return `
	<div class="border rounded-lg p-6 hover:shadow-md transition-shadow">
		<div class="flex justify-between items-start">
			<div>
				<p class="text-lg font-[Playfair_Display] mb-2">${escapedString(quote.text)}</p>
				<p class="text-sm text-gray-800">${escapedString(quote.author)}</p>
			</div>
			<div class="flex gap-2">
				<button class="text-gray-600 hover:text-custom" onclick="quoteEdit('${quote.id}');">
					<i class="fas fa-edit"></i>
				</button>
				<button class="text-gray-600 hover:text-red-500" onclick="quoteDelete('${quote.id}');">
					<i class="fas fa-trash"></i>
				</button>
			</div>
		</div>
		<div class="w-full flex justify-between items-center flex-wrap gap-1 mt-2">
			<p class="text-xs text-gray-400">${new Date(quote.datetime).toLocaleString()}</p>
			<p class="text-xs text-gray-400">Добавил ${quote.creator.name}</p>
		</div>
	</div>
	`;
}

function load() {
	var quotesCount;
	$("#input-search").val(search);
	$("#input-sorting").val(sorting);
	container = $("#block-quotes");
	$.ajax({
		async: false,
		url: `/api/quotes?filter=${encodeURIComponent(search)}&sort=${encodeURIComponent(sorting)}&limit=${PAGE_SIZE}&offset=${(currPage - 1)*PAGE_SIZE}`,
		type: "GET",
		dataType: "json",
		success: function (resp) {
			quotesCount = resp.pagination.totalCount
			if (resp.pagination.count == 0) {
				container.html("<p style='text-align: center;'><i>Чёт нету ничего...</i></p>");
				return;
			}
			resp.quotes.forEach((quote) => {
				container.append(renderBlockQuote(quote));
			});
		},
		error: function (err) {
			$("#error-message").text(err.responseJSON.error);
			$("#error").removeClass("hidden");
		},
		complete: function () {
			$("#block-quotes-loader").addClass("hidden");
		},
	});
	totalPages = Math.ceil(quotesCount / PAGE_SIZE);
	$("#btn-page-curr").text(currPage);
	if (currPage > 1) {
		$("#btn-page-first").removeClass("hidden");
		if (currPage > 2) {
			$("#btn-page-prev").text(currPage - 1);
			$("#btn-page-prev").removeClass("hidden");
			if (currPage > 3) {
				$("#pages-prev").removeClass("hidden");
			}
		}
	}
	if (currPage < totalPages) {
		$("#btn-page-last").text(totalPages);
		$("#btn-page-last").removeClass("hidden");
		if (currPage < totalPages - 1) {
			$("#btn-page-next").text(currPage + 1);
			$("#btn-page-next").removeClass("hidden");
			if (currPage < totalPages - 2) {
				$("#pages-next").removeClass("hidden");
			}
		}
	}
}

function reload() {
	container = $("#block-quotes");
	loader = $("#block-quotes-loader");
	loader.removeClass("hidden");
	container.html(loader);
	$("#error").addClass("hidden");
	$("#btn-page-first").addClass("hidden");
	$("#pages-prev").addClass("hidden");
	$("#btn-page-prev").addClass("hidden");
	$("#btn-page-curr").text(1);
	$("#btn-page-next").addClass("hidden");
	$("#pages-next").addClass("hidden");
	$("#btn-page-last").addClass("hidden");
	load();
}

function quoteEdit(quote_id) {
	$.ajax({
		url: `/api/quotes/${quote_id}`,
		type: "GET",
		dataType: "json",
		success: function (resp) {
			temp_quote_id = quote_id;
			$("#edit-quote-text").val(resp.text);
			$("#edit-quote-author").val(resp.author);
			$("#edit-quote-datetime").val(resp.datetime.slice(0,19));
			$("body").css("overflow", "hidden");
			$("#quote-editor").css("top", $(window).scrollTop());
			$("#quote-editor").removeClass("hidden");
		},
		error: function (err) {
			$("#quote-editor-error-message").text(err.responseJSON.error);
			$("#quote-editor-error").removeClass("hidden");
		},
	});
}

function quoteDelete(quote_id) {
	$.ajax({
		url: `/api/quotes/${quote_id}`,
		type: "DELETE",
		success: function (resp) {
			reload();
			$("#error").addClass("hidden");
		},
		error: function (err) {
			$("#error-message").text(err.responseJSON.error);
			$("#error").removeClass("hidden");
		},
	});
}

$(document).on("click", "#btn-add-open", function (e) {
	now = new Date;
	now = new Date(now.getTime() - now.getTimezoneOffset() * 60000);
	$("#new-quote-datetime").val(now.toJSON().slice(0,19));
	$("body").css("overflow", "hidden");
	$("#quote-creator").css("top", $(window).scrollTop());
	$("#quote-creator").removeClass("hidden");
});

$(document).on("click", "#btn-add-close", function (e) {
	$("#quote-creator").addClass("hidden");
	$("body").css("overflow", "");
});

$(document).on("submit", "#quote-create", function (e) {
	e.preventDefault();
	data = formToJSON($("#quote-create"));
	data.datetime = datetimeToLocalISO(data.datetime);
	$.ajax({
		url: "/api/quotes",
		type: "POST",
		contentType: "application/json",
		data: JSON.stringify(data),
		processData: false,
		dataType: "json",
		success: function (resp) {
			$("#quote-creator").addClass("hidden");
			$("body").css("overflow", "");
			reload();
			$("#new-quote-text").val("");
			$("#new-quote-author").val("");
		},
		error: function (err) {
			$("#quote-creator-error-message").text(err.responseJSON.error);
			$("#quote-creator-error").removeClass("hidden");
		},
	});
});

$(document).on("submit", "#quote-update", function (e) {
	e.preventDefault();
	data = formToJSON($("#quote-update"));
	data.datetime = datetimeToLocalISO(data.datetime);
	$.ajax({
		url: `/api/quotes/${temp_quote_id}`,
		type: "PATCH",
		contentType: "application/json",
		data: JSON.stringify(data),
		processData: false,
		dataType: "json",
		success: function (resp) {
			$("#quote-editor").addClass("hidden");
			$("body").css("overflow", "");
			reload();
			$("#new-quote-text").val("");
			$("#new-quote-author").val("");
		},
		error: function (err) {
			$("#quote-editor-error-message").text(err.responseJSON.error);
			$("#quote-editor-error").removeClass("hidden");
		},
	});
});

$(document).on("click", "#btn-edit-close", function (e) {
	$("body").css("overflow", "");
	$("#quote-editor").addClass("hidden");
	$("#quote-editor textarea,input").val("");
	$("#quote-editor-error").addClass("hidden");
});

$(window).on("load", function (e) {
	load();
});

$(document).on("click", "#btn-refresh", function (e) {
	search = $("#input-search").val();
	if (search != "") {
		currPage = 1;
		sessionStorage.setItem("search", currPage);
	}
	sorting = $("#input-sorting option:selected").val();
	reload();
	sessionStorage.setItem("search", search);
	sessionStorage.setItem("sort", sorting);
});

$(document).on("click", "#btn-page-first", function (e) {
	currPage = 1;
	reload();
	sessionStorage.setItem("page", currPage);
});

$(document).on("click", "#btn-page-prev", function (e) {
	currPage--;
	reload();
	sessionStorage.setItem("page", currPage);
});

$(document).on("click", "#btn-page-next", function (e) {
	currPage++;
	reload();
	sessionStorage.setItem("page", currPage);
});

$(document).on("click", "#btn-page-last", function (e) {
	currPage = totalPages;
	reload();
	sessionStorage.setItem("page", currPage);
});
