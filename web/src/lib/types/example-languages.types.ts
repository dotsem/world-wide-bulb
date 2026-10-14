export const Language = {
	Python: 'python',
	JavaScript: 'javascript',
	Go: 'go',
	Curl: 'curl'
} as const;

export type Language = (typeof Language)[keyof typeof Language];
