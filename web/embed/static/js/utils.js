function datetimeToLocalISO(datetime) {
	options = {
		year: "numeric",
		month: "2-digit",
		day: "2-digit",
		hour: "2-digit",
		minute: "2-digit",
		second: "2-digit",
		timeZoneName: "longOffset",
	};
	formatter = new Intl.DateTimeFormat("iso", options);
	date = new Date(datetime);
	parts = {}
	formatter
		.formatToParts(date)
		.map(({ type, value }) => {
			if (type === "timeZoneName") {
				value = value.slice(3);
			}
			if (type !== "literal") {
				parts[type] = value;
			}
		});
	return `${parts.year}-${parts.month}-${parts.day}T${parts.hour}:${parts.minute}:${parts.second}${parts.timeZoneName}`;
}

function escapedString(str) {
	return str
		.replaceAll("&", "&amp;")
		.replaceAll("<", "&lt;")
		.replaceAll(">", "&gt;")
		.replaceAll("\n", "<br>");
}

function formToJSON(form) {
	formdata = form.serializeArray();
	data = {};
	$(formdata).each(function (index, obj) {
		data[obj.name] = obj.value;
	});
    return data;
}
