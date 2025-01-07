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

function escapedString(str) {
	return str
		.replace("&", "&amp;")
		.replace("<", "&lt;")
		.replace(">", "&gt;")
		.replace("\n", "<br>");
}

function renderBlockQuote(quote) {
	return `
	<div class="border rounded-lg p-6 hover:shadow-md transition-shadow" quote_id="${quote.id}">
		<div class="flex justify-between items-start">
			<div>
				<p class="text-lg font-[Playfair_Display] mb-2">${escapedString(quote.text)}</p>
				<p class="text-sm text-gray-800">${escapedString(quote.author)}</p>
				<p class="text-xs text-gray-400 mt-2">${quote.datetime}</p>
			</div>
			<div class="flex gap-2">
				<button class="text-gray-600 hover:text-custom">
					<i class="fas fa-edit"></i>
				</button>
				<button class="text-gray-600 hover:text-red-500">
					<i class="fas fa-trash"></i>
				</button>
			</div>
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
