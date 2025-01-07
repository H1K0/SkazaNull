function datetimeToLocalISO(datetime) {
	var options = {
		year: "numeric",
		month: "2-digit",
		day: "2-digit",
		hour: "2-digit",
		minute: "2-digit",
		second: "2-digit",
		timeZoneName: "longOffset",
	};
	var formatter = new Intl.DateTimeFormat("sv-SE", options);
	var date = new Date(datetime);
	return formatter
		.formatToParts(date)
		.map(({ type, value }) => {
			if (type === "timeZoneName") {
				return value.slice(3);
			} else {
				return value;
			}
		})
		.join("")
		.replace(" ", "T")
		.replace(" ", "");
}

function escapedString(str) {
	return str
		.replace("&", "&amp;")
		.replace("<", "&lt;")
		.replace(">", "&gt;")
		.replace("\n", "<br>");
}

function formToJSON(form) {
	formdata = form.serializeArray();
	data = {};
	$(formdata).each(function (index, obj) {
		data[obj.name] = obj.value;
	});
    return data;
}
