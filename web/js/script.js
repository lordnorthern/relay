let socket = new WebSocket("ws://www.elkz.net:7717/ws")
socket.onopen = () => {
	logMessage("Successfully Connected");
};

socket.onclose = event => {
	logMessage(event);
	socket.send("Client Closed!")
};

socket.onerror = error => {
	logMessage(error);
};

socket.onmessage = event => {
	logMessage(event.data);
};

function logMessage(message) {
	let nClass = "json"
	let actualResponse = ""
	let headers = []
	try {
		const content = JSON.parse(message);
		try {
			actualResponse = JSON.parse(content.body);
			
			try {
				headers = JSON.parse(content.header)
			} catch {}

		} catch (error) {
			nClass = "text"
			actualResponse = message
		}
	} catch (error) {
		nClass = "text"
		actualResponse = message
	}

	const date = new moment;
	const randomNumber = Math.floor(Math.random() * 999999999999999);
	message = `<div>------------------  ${date.format("YYYY-MM-DD HH-mm-ss")}  ------------------</div>
		<pre class="actual-response">${actualResponse}</pre>
		<div>--------------------------------------------</div>`
	$("#responses").prepend(`<li class="${nClass} li-${randomNumber}">${message}</li>`);
	if (nClass === "json") {
		$(`.li-${randomNumber}.json .actual-response`).jsonBrowse({
			header: headers,
			body: actualResponse,
		}, { collapsed: false, withQuotes: true });
	}
}