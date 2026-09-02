// Connect to World Wide Bulb SSE stream in Node.js (v22.3+ with --experimental-eventsource) or Browser
const sseUrl = process.env.WWB_URL || 'https://wwb.dotsem.be/api/v1/events';

console.log(`Connecting to ${sseUrl}...`);
const events = new EventSource(sseUrl);

events.addEventListener('state_changed', (e) => {
	const data = JSON.parse(e.data);
	console.log('Lamp state changed:', data.state ? 'ON' : 'OFF', data);
});

events.addEventListener('reason_updated', (e) => {
	const data = JSON.parse(e.data);
	console.log('Reason updated:', data.reason);
});

events.onerror = (err) => {
	console.error('SSE connection error:', err);
};
