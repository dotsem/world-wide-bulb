const baseUrl = (process.env.WWB_URL || 'http://localhost:8080')
	.replace(/\/api\/v1(?:\/events)?\/?$/, '')
	.replace(/\/+$/, '');

// region: get_state
async function getState() {
	const res = await fetch(`${baseUrl}/api/v1/state`);
	if (!res.ok) throw new Error(`failed to get state: ${res.status}`);
	const data = await res.json();
	console.log(`State: ${data.state ? 'ON' : 'OFF'} (${data.viewers} viewers)`);
	return data;
}
// endregion: get_state

// region: get_history
async function getHistory(limit = 10) {
	const res = await fetch(`${baseUrl}/api/v1/history?limit=${limit}`);
	if (!res.ok) throw new Error(`failed to get history: ${res.status}`);
	const data = await res.json();
	console.log(`History entries: ${data.items ? data.items.length : 0}`);
	return data;
}
// endregion: get_history

// region: post_toggle
async function toggleBulb() {
	const res = await fetch(`${baseUrl}/api/v1/toggle`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' }
	});
	if (!res.ok) throw new Error(`failed to toggle bulb: ${res.status}`);
	const data = await res.json();
	console.log(`Toggled bulb, ID: ${data.id} (State: ${data.state ? 'ON' : 'OFF'})`);
	return data;
}
// endregion: post_toggle

// region: post_reason
async function submitReason(toggleId, reason) {
	const res = await fetch(`${baseUrl}/api/v1/reason`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ id: toggleId, reason })
	});
	if (!res.ok) throw new Error(`failed to submit reason: ${res.status}`);
	const data = await res.json();
	console.log(`Reason updated for ${toggleId}`);
	return data;
}
// endregion: post_reason

async function main() {
	console.log(`Connecting to REST API at: ${baseUrl}`);
	await getState();
	await getHistory();
	const toggle = await toggleBulb();
	if (toggle.id) {
		await submitReason(toggle.id, 'Turning on the bulb for client demo');
	}
}

main().catch((err) => {
	console.error('Error in REST client:', err.message);
	process.exit(1);
});
