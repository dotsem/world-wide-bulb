import { codeToHtml } from 'shiki';
import jsCode from '../../../../examples/sse/javascript/client.js?raw';
import pyCode from '../../../../examples/sse/python/client.py?raw';
import goCode from '../../../../examples/sse/go/main.go?raw';

export const prerender = true;

export async function load() {
	const theme = 'tokyo-night';

	const [jsHtml, pyHtml, goHtml, curlHtml, stateChangedHtml, reasonUpdatedHtml, pingHtml] =
		await Promise.all([
			codeToHtml(jsCode, {
				lang: 'javascript',
				theme,
				defaultColor: false
			}),
			codeToHtml(pyCode, {
				lang: 'python',
				theme,
				defaultColor: false
			}),
			codeToHtml(goCode, {
				lang: 'go',
				theme,
				defaultColor: false
			}),
			codeToHtml('curl -N https://wwb.dotsem.be/api/v1/events', {
				lang: 'bash',
				theme,
				defaultColor: false
			}),
			codeToHtml(
				JSON.stringify(
					{
						type: 'state_changed',
						id: 42,
						state: true,
						reason: 'It was getting dark here',
						created_at: '2026-09-02T17:35:00Z'
					},
					null,
					2
				),
				{
					lang: 'json',
					theme,
					defaultColor: false
				}
			),
			codeToHtml(
				JSON.stringify(
					{
						type: 'reason_updated',
						id: 42,
						reason: 'Party at the office!'
					},
					null,
					2
				),
				{
					lang: 'json',
					theme,
					defaultColor: false
				}
			),
			codeToHtml('""', {
				lang: 'json',
				theme,
				defaultColor: false
			})
		]);

	return {
		raw: {
			javascript: jsCode,
			python: pyCode,
			go: goCode,
			curl: 'curl -N https://wwb.dotsem.be/api/v1/events'
		},
		highlighted: {
			javascript: jsHtml,
			python: pyHtml,
			go: goHtml,
			curl: curlHtml
		},
		payloads: {
			stateChanged: stateChangedHtml,
			reasonUpdated: reasonUpdatedHtml,
			ping: pingHtml
		}
	};
}
