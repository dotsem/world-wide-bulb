<script lang="ts">
	let {
		isOn = false,
		disabled = false,
		durationMs = 500,
		class: className = 'h-56 w-56 sm:h-72 sm:w-72 md:h-80 md:w-80 lg:h-96 lg:w-96',
		onclick
	}: {
		isOn?: boolean;
		disabled?: boolean;
		durationMs?: number;
		class?: string;
		onclick?: () => void;
	} = $props();
</script>

<button
	type="button"
	{onclick}
	aria-disabled={disabled}
	style="--duration: {durationMs}ms;"
	class="group relative inline-flex items-center justify-center rounded-full p-8 transition-all focus:outline-none duration-(--duration) {disabled
		? 'cursor-not-allowed'
		: 'cursor-pointer'} {className}"
	aria-label={isOn ? 'Turn off lightbulb' : 'Turn on lightbulb'}
>
	<!-- Outer Broad Ambient Glow -->
	<div
		style="transition-duration: var(--duration);"
		class="absolute inset-0 rounded-full bg-amber-400/25 blur-3xl transition-all duration-300 pointer-events-none"
		class:opacity-100={isOn}
		class:opacity-0={!isOn}
		class:scale-150={isOn}
		class:scale-75={!isOn}
	></div>

	<!-- Inner Core Ambient Glow -->
	<div
		style="transition-duration: var(--duration);"
		class="absolute inset-4 rounded-full bg-amber-300/35 blur-xl transition-all duration-300 pointer-events-none"
		class:opacity-100={isOn}
		class:opacity-0={!isOn}
		class:scale-110={isOn}
		class:scale-90={!isOn}
	></div>

	<!-- Vector Lightbulb SVG -->
	<svg
		version="1.1"
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 600 600"
		style="transition-duration: var(--duration);"
		class="relative h-full w-full transition-all transform group-hover:scale-105 group-active:scale-95"
		class:drop-shadow-[0_0_40px_rgba(251,191,36,0.85)]={isOn}
		class:drop-shadow-none={!isOn}
	>
		<g id="objects">
			<g>
				<!-- Inner Filament Glass Mount Stem -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-all stroke-[2px] {isOn
						? 'fill-amber-400/30 stroke-amber-300/60'
						: 'fill-slate-400/15 stroke-slate-400/30'}"
					d="M 285,418 L 285,345 C 285,335 315,335 315,345 L 315,418 Z"
				/>

				<!-- Inner Filament Support Wires & Coiled Filament Loop -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-all fill-none {isOn
						? 'stroke-[#FFFBEB] stroke-[3.5px] drop-shadow-[0_0_10px_rgba(255,251,235,0.9)]'
						: 'stroke-slate-300/40 stroke-[2px]'}"
					d="M 290,345 L 272,215 M 310,345 L 328,215"
				/>
				<path
					style="transition-duration: var(--duration);"
					class="transition-all fill-none {isOn
						? 'stroke-[#FFFBEB] stroke-[4.5px] drop-shadow-[0_0_14px_rgba(255,251,235,1)]'
						: 'stroke-slate-300/60 stroke-[2.5px]'}"
					d="M 272,215 C 272,170 328,170 328,215"
				/>

				<!-- Main Bulb Glass Body (Translucent when OFF, Glowing Yellow when ON) -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-all {isOn
						? 'fill-[#FDC10D] stroke-amber-300/40 stroke-[3px]'
						: 'fill-slate-100/8 stroke-slate-300/50 stroke-[4px]'}"
					d="M435.995,186.014c0.22,30.926-9.883,59.482-27.068,82.433
					c-21.043,28.103-32.929,61.996-32.929,97.104v4.452c0,16.569-13.431,30-30,30h-92c-16.569,0-30-13.431-30-30v-2.959
					c0-35.773-11.975-70.352-33.29-99.082c-16.783-22.622-26.71-50.631-26.71-80.96c0-75.441,61.428-136.537,136.995-135.996
					C375.022,51.536,435.469,111.985,435.995,186.014z"
				/>

				<!-- Screw Base Rings -->
				<path
					style="fill:#575553;"
					d="M355.998,418.003h-112c-3.314,0-6-2.686-6-6v-6c0-3.314,2.686-6,6-6h112c3.314,0,6,2.686,6,6v6C361.998,415.316,359.312,418.003,355.998,418.003z"
				/>
				<path
					style="fill:#575553;"
					d="M355.998,436.003h-112c-3.314,0-6-2.686-6-6v-6c0-3.314,2.686-6,6-6h112c3.314,0,6,2.686,6,6v6C361.998,433.316,359.312,436.003,355.998,436.003z"
				/>
				<path
					style="fill:#575553;"
					d="M355.998,454.003h-112c-3.314,0-6-2.686-6-6v-6c0-3.314,2.686-6,6-6h112c3.314,0,6,2.686,6,6v6C361.998,451.316,359.312,454.003,355.998,454.003z"
				/>
				<path
					style="fill:#575553;"
					d="M355.998,472.003h-112c-3.314,0-6-2.686-6-6v-6c0-3.314,2.686-6,6-6h112c3.314,0,6,2.686,6,6v6C361.998,469.316,359.312,472.003,355.998,472.003z"
				/>
				<path
					style="fill:#575553;"
					d="M335.851,490.003h-71.706c-5.646,0-10.885-2.68-13.829-7.075l-7.318-10.925h114l-7.318,10.925C346.736,487.322,341.497,490.003,335.851,490.003z"
				/>
				<path
					style="fill:#575553;"
					d="M313.16,507.003h-26.324c-6.524,0-12.603-2.94-16.147-7.808l-6.691-9.192h72l-6.691,9.192C325.763,504.063,319.684,507.003,313.16,507.003z"
				/>

				<!-- Right Side Rim Reflection -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-all {isOn ? 'fill-[#EAA60F] opacity-100' : 'fill-white/12 opacity-80'}"
					d="M435.998,186.013c0.22,30.93-9.89,59.48-27.07,82.43c-21.04,28.11-32.93,62-32.93,97.11v4.45c0,16.57-13.43,30-30,30h-92c-16.57,0-30-13.43-30-30v-2.96c0-1.53-0.02-3.06-0.08-4.59c5.79,7.36,15.31,15.55,29.08,15.55c24.87,0,48.17-8.37,66.86-22.74c18.68-14.38,32.76-34.76,39.16-58.8c4.46-16.76,11.737-33.687,18.98-50.46c38-88,8-206-135.73-182.17c17.8-8.37,37.72-12.98,58.72-12.83C375.018,51.533,435.468,111.983,435.998,186.013z"
				/>

				<!-- Top-Left Glossy Glass Glare / Specular Highlight Paths -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-all {isOn ? 'fill-[#FFE67D] opacity-100' : 'fill-white/30 opacity-90'}"
					d="M242,285c0,0-100.101-55.825-42.511-159.922c1.906-3.446,6.164-4.826,9.724-3.144l22.259,10.511
					c3.875,1.83,5.408,6.53,3.371,10.3C225.474,160.092,206.19,210.197,242,285z"
				/>
				<path
					style="transition-duration: var(--duration);"
					class="transition-all {isOn ? 'fill-[#FFE67D] opacity-100' : 'fill-white/45 opacity-95'}"
					d="M215.412,108.706l22.549,11.274c2.292,1.146,5.062,0.543,6.676-1.447
					c3.118-3.847,8.937-10.334,17.103-16.539c3.179-2.416,2.728-7.333-0.82-9.164l-20.567-10.615
					c-1.549-0.799-3.379-0.809-4.939-0.032c-4.599,2.291-14.425,8.047-21.968,18.422
					C211.461,103.334,212.394,107.197,215.412,108.706z"
				/>
			</g>
		</g>
	</svg>
</button>
