let socket = new WebSocket("wss://www.elkz.net:7717/ws")
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
	try {
		message = JSON.stringify(JSON.parse(message), null, 2);
	} catch (error) {}
	const date = new moment;
	message = "\n------------------" + date.format("YYYY-MM-DD HH-mm-ss") + "------------------\n" +
		message +
		"\n--------------------------------------------\n"
	$("#console").prepend(message + "\n");
}