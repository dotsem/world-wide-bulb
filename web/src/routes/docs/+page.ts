import { codeToHtml } from 'shiki';
import sseJsCode from '../../../../examples/sse/javascript/client.js?raw';
import ssePyCode from '../../../../examples/sse/python/client.py?raw';
import sseGoCode from '../../../../examples/sse/go/main.go?raw';
import restJsCode from '../../../../examples/rest/javascript/client.js?raw';
import restPyCode from '../../../../examples/rest/python/client.py?raw';
import restGoCode from '../../../../examples/rest/go/main.go?raw';

export const prerender = true;

function extractRegion(code: string, region: string): string {
	const regex = new RegExp(
		`(?:\\/\\/|#)\\s*region:\\s*${region}\\r?\\n([\\s\\S]*?)(?:\\/\\/|#)\\s*endregion:\\s*${region}`,
		'i'
	);
	const match = code.match(regex);
	return match ? match[1].trim() : '';
}

export async function load() {
	const theme = 'tokyo-night';

	async function highlight(code: string, lang: string) {
		return codeToHtml(code, { lang, theme, defaultColor: false });
	}

	async function highlightJson(obj: unknown) {
		return highlight(JSON.stringify(obj, null, 2), 'json');
	}

	async function makeExamples(curlCmd: string, regionKey: string, repoBase: 'rest' | 'sse') {
		const rawCurl = curlCmd.trim();
		const rawJs = repoBase === 'rest' ? extractRegion(restJsCode, regionKey) : sseJsCode;
		const rawPy = repoBase === 'rest' ? extractRegion(restPyCode, regionKey) : ssePyCode;
		const rawGo = repoBase === 'rest' ? extractRegion(restGoCode, regionKey) : sseGoCode;

		const jsPath =
			repoBase === 'rest'
				? 'examples/rest/javascript/client.js'
				: 'examples/sse/javascript/client.js';
		const pyPath =
			repoBase === 'rest' ? 'examples/rest/python/client.py' : 'examples/sse/python/client.py';
		const goPath = repoBase === 'rest' ? 'examples/rest/go/main.go' : 'examples/sse/go/main.go';

		const [curlHtml, jsHtml, pyHtml, goHtml] = await Promise.all([
			highlight(rawCurl, 'bash'),
			highlight(rawJs, 'javascript'),
			highlight(rawPy, 'python'),
			highlight(rawGo, 'go')
		]);

		return {
			curl: { highlighted: curlHtml, raw: rawCurl },
			javascript: { highlighted: jsHtml, raw: rawJs, repoPath: jsPath },
			python: { highlighted: pyHtml, raw: rawPy, repoPath: pyPath },
			go: { highlighted: goHtml, raw: rawGo, repoPath: goPath }
		};
	}

	const wsRawJs = `const ws = new WebSocket('wss://wwb.dotsem.be/ws');

ws.onopen = () => console.log('Connected to WebSocket stream');
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message.type, message);
};
ws.onclose = () => console.log('Disconnected');`;

	const wsRawPy = `import asyncio
import json
import websockets

async def listen():
    uri = "wss://wwb.dotsem.be/ws"
    async with websockets.connect(uri) as ws:
        print("Connected to WebSocket stream")
        async for msg in ws:
            data = json.loads(msg)
            print("Received:", data.get("type"), data)

asyncio.run(listen())`;

	const wsRawGo = `package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/gorilla/websocket"
)

func main() {
	url := "wss://wwb.dotsem.be/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		var payload map[string]any
		_ = json.Unmarshal(msg, &payload)
		fmt.Printf("Received: %v\\n", payload)
	}
}`;

	const [
		stateExamples,
		historyExamples,
		toggleExamples,
		reasonExamples,
		sseExamples,
		stateRespJson,
		historyRespJson,
		toggleRespJson,
		reasonRespJson,
		sseStatePayload,
		sseReasonPayload,
		ssePingPayload,
		wsCurlHtml,
		wsJsHtml,
		wsPyHtml,
		wsGoHtml,
		wsStatePayload,
		wsReasonPayload
	] = await Promise.all([
		makeExamples('curl -s https://wwb.dotsem.be/api/v1/state', 'get_state', 'rest'),
		makeExamples('curl -s "https://wwb.dotsem.be/api/v1/history?limit=10"', 'get_history', 'rest'),
		makeExamples(
			'curl -X POST https://wwb.dotsem.be/api/v1/toggle \\\n  -H "Content-Type: application/json"',
			'post_toggle',
			'rest'
		),
		makeExamples(
			'curl -X POST https://wwb.dotsem.be/api/v1/reason \\\n  -H "Content-Type: application/json" \\\n  -d \'{"id": "a62c1726-8da9-4ab1-8a92-8e61533c0013", "reason": "Office party"}\'',
			'post_reason',
			'rest'
		),
		makeExamples('curl -N https://wwb.dotsem.be/api/v1/events', '', 'sse'),
		highlightJson({
			state: true,
			cooldown_ms: 0,
			viewers: 24
		}),
		highlightJson({
			items: [
				{
					id: 42,
					state: true,
					reason: 'It was getting dark here',
					created_at: '2026-10-14T08:00:00Z'
				}
			],
			next_before: 22
		}),
		highlightJson({
			id: 'a62c1726-8da9-4ab1-8a92-8e61533c0013',
			state: true,
			created_at: '2026-10-14T08:00:00Z',
			cooldown_ms: 10000
		}),
		highlightJson({
			status: 'ok'
		}),
		highlightJson({
			type: 'state_changed',
			id: 42,
			state: true,
			reason: 'It was getting dark here',
			created_at: '2026-10-14T08:00:00Z'
		}),
		highlightJson({
			type: 'reason_updated',
			id: 42,
			reason: 'Party at the office!'
		}),
		highlightJson(''),
		highlight('wscat -c wss://wwb.dotsem.be/ws', 'bash'),
		highlight(wsRawJs, 'javascript'),
		highlight(wsRawPy, 'python'),
		highlight(wsRawGo, 'go'),
		highlightJson({
			type: 'state_changed',
			id: 42,
			state: true,
			reason: 'It was getting dark here',
			created_at: '2026-10-14T08:00:00Z'
		}),
		highlightJson({
			type: 'reason_updated',
			id: 42,
			reason: 'Party at the office!'
		})
	]);

	const wsExamples = {
		curl: { highlighted: wsCurlHtml, raw: 'wscat -c wss://wwb.dotsem.be/ws' },
		javascript: { highlighted: wsJsHtml, raw: wsRawJs },
		python: { highlighted: wsPyHtml, raw: wsRawPy },
		go: { highlighted: wsGoHtml, raw: wsRawGo }
	};

	return {
		endpoints: {
			state: {
				examples: stateExamples,
				responseHtml: stateRespJson
			},
			history: {
				examples: historyExamples,
				responseHtml: historyRespJson
			},
			toggle: {
				examples: toggleExamples,
				responseHtml: toggleRespJson
			},
			reason: {
				examples: reasonExamples,
				responseHtml: reasonRespJson
			},
			events: {
				examples: sseExamples,
				responseHtml: sseStatePayload,
				payloads: {
					stateChanged: sseStatePayload,
					reasonUpdated: sseReasonPayload,
					ping: ssePingPayload
				}
			},
			ws: {
				examples: wsExamples,
				responseHtml: wsStatePayload,
				payloads: {
					stateChanged: wsStatePayload,
					reasonUpdated: wsReasonPayload
				}
			}
		}
	};
}
