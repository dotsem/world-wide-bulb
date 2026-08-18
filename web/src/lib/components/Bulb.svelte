<script lang="ts">
	let {
		isOn = false,
		disabled = false,
		durationMs = 200,
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
	class="group relative inline-flex items-center justify-center rounded-full p-8 focus:outline-none touch-manipulation {disabled
		? 'cursor-not-allowed'
		: 'cursor-pointer'} {className}"
	aria-label={isOn ? 'Turn off lightbulb' : 'Turn on lightbulb'}
>
	<!-- Outer Broad Ambient Glow -->
	<div
		style="transition-duration: var(--duration);"
		class="absolute inset-0 rounded-full bg-[radial-gradient(circle,rgba(251,191,36,0.35)_0%,transparent_70%)] transition pointer-events-none will-change-transform"
		class:opacity-100={isOn}
		class:opacity-0={!isOn}
		class:scale-150={isOn}
		class:scale-75={!isOn}
	></div>

	<!-- Inner Core Ambient Glow -->
	<div
		style="transition-duration: var(--duration);"
		class="absolute inset-4 rounded-full bg-[radial-gradient(circle,rgba(251,191,36,0.5)_0%,transparent_60%)] transition pointer-events-none will-change-transform"
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
		class="relative h-full w-full transition-transform duration-200 transform group-hover:scale-105 group-active:scale-95"
		class:md:drop-shadow-[0_0_40px_rgba(251,191,36,0.85)]={isOn}
		class:drop-shadow-none={!isOn}
	>
		<g id="objects">
			<g>
				<!-- Inner Filament Glass Mount Stem -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-colors stroke-[2px] {isOn
						? 'fill-app-accent/30 stroke-app-accent/60'
						: 'fill-app-muted/15 stroke-app-muted/30'}"
					d="M 285,418 L 285,345 C 285,335 315,335 315,345 L 315,418 Z"
				/>

				<!-- Inner Filament Support Wires & Coiled Filament Loop -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-colors fill-none {isOn
						? 'stroke-[#FFFBEB] stroke-[3.5px] md:drop-shadow-[0_0_10px_rgba(255,251,235,0.9)]'
						: 'stroke-app-muted/40 stroke-[2px]'}"
					d="M 290,345 L 272,215 M 310,345 L 328,215"
				/>
				<path
					style="transition-duration: var(--duration);"
					class="transition-colors fill-none {isOn
						? 'stroke-[#FFFBEB] stroke-[4.5px] md:drop-shadow-[0_0_14px_rgba(255,251,235,1)]'
						: 'stroke-app-muted/60 stroke-[2.5px]'}"
					d="M 272,215 C 272,170 328,170 328,215"
				/>

				<!-- Main Bulb Glass Body (Translucent when OFF, Glowing Yellow when ON) -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-colors {isOn
						? 'fill-app-bulb-fill stroke-app-accent/40 stroke-[3px]'
						: 'fill-app-surface-hover/30 stroke-app-muted/50 stroke-[4px]'}"
					d="M435.995,186.014c0.22,30.926-9.883,59.482-27.068,82.433
					c-21.043,28.103-32.929,61.996-32.929,97.104v4.452c0,16.569-13.431,30-30,30h-92c-16.569,0-30-13.431-30-30v-2.959
					c0-35.773-11.975-70.352-33.29-99.082c-16.783-22.622-26.71-50.631-26.71-80.96c0-75.441,61.428-136.537,136.995-135.996
					C375.022,51.536,435.469,111.985,435.995,186.014z"
				/>

				<!-- Screw Base Rings -->
				<path
					class="fill-app-bulb-base"
					d="M355.998,418.003h-112c-3.314,0-6-2.686-6-6v-6c0-3.314,2.686-6,6-6h112c3.314,0,6,2.686,6,6v6C361.998,415.316,359.312,418.003,355.998,418.003z"
				/>
				<path
					class="fill-app-bulb-base"
					d="M355.998,436.003h-112c-3.314,0-6-2.686-6-6v-6c0-3.314,2.686-6,6-6h112c3.314,0,6,2.686,6,6v6C361.998,433.316,359.312,436.003,355.998,436.003z"
				/>
				<path
					class="fill-app-bulb-base"
					d="M355.998,454.003h-112c-3.314,0-6-2.686-6-6v-6c0-3.314,2.686-6,6-6h112c3.314,0,6,2.686,6,6v6C361.998,451.316,359.312,454.003,355.998,454.003z"
				/>
				<path
					class="fill-app-bulb-base"
					d="M355.998,472.003h-112c-3.314,0-6-2.686-6-6v-6c0-3.314,2.686-6,6-6h112c3.314,0,6,2.686,6,6v6C361.998,469.316,359.312,472.003,355.998,472.003z"
				/>
				<path
					class="fill-app-bulb-base"
					d="M335.851,490.003h-71.706c-5.646,0-10.885-2.68-13.829-7.075l-7.318-10.925h114l-7.318,10.925C346.736,487.322,341.497,490.003,335.851,490.003z"
				/>
				<path
					class="fill-app-bulb-base"
					d="M313.16,507.003h-26.324c-6.524,0-12.603-2.94-16.147-7.808l-6.691-9.192h72l-6.691,9.192C325.763,504.063,319.684,507.003,313.16,507.003z"
				/>

				<!-- Screw Base Ring Bevel Shadows -->
				<path
					class="fill-app-bulb-base-dark"
					d="M361.998,406.003v6c0,3.31-2.69,6-6,6h-112c-3.31,0-6-2.69-6-6v-0.81c0.88,0.52,1.9,0.81,3,0.81h112c3.31,0,6-2.69,6-6v-5.19C360.788,401.843,361.998,403.783,361.998,406.003z"
				/>
				<path
					class="fill-app-bulb-base-dark"
					d="M361.998,424.003v6c0,3.31-2.69,6-6,6h-112c-3.31,0-6-2.69-6-6v-0.81c0.88,0.52,1.9,0.81,3,0.81h112c3.31,0,6-2.69,6-6v-5.19C360.788,419.843,361.998,421.783,361.998,424.003z"
				/>
				<path
					class="fill-app-bulb-base-dark"
					d="M361.998,442.003v6c0,3.31-2.69,6-6,6h-112c-3.31,0-6-2.69-6-6v-0.81c0.88,0.52,1.9,0.81,3,0.81h112c3.31,0,6-2.69,6-6v-5.19C360.788,437.843,361.998,439.783,361.998,442.003z"
				/>
				<path
					class="fill-app-bulb-base-dark"
					d="M361.998,460.003v6c0,3.31-2.69,6-6,6h-112c-3.31,0-6-2.69-6-6v-0.81c0.88,0.52,1.9,0.81,3,0.81h112c3.31,0,6-2.69,6-6v-5.19C360.788,455.843,361.998,457.783,361.998,460.003z"
				/>
				<path
					class="fill-app-bulb-base-dark"
					d="M356.998,472.003l-7.32,10.92c-2.94,4.4-8.18,7.08-13.83,7.08h-71.7c-5.65,0-10.89-2.68-13.83-7.08l-0.1-0.14c2.07,0.79,4.33,1.22,6.67,1.22h74.22c5.85,0,11.27-2.68,14.32-7.08l3.41-4.92H356.998z"
				/>
				<path
					class="fill-app-bulb-base-dark"
					d="M335.998,490.003l-6.69,9.19c-3.55,4.87-9.62,7.81-16.15,7.81h-26.32c-4.94,0-9.61-1.68-13.13-4.61c2.5,1.04,5.28,1.61,8.13,1.61h26.32c6.53,0,12.6-2.94,16.15-7.81l4.51-6.19H335.998z"
				/>

				<!-- Screw Base Highlights -->
				<path
					class="fill-app-bulb-base-light"
					d="M264.998,407.003h-21c-1.105,0-2-0.895-2-2l0,0c0-1.105,0.895-2,2-2h21c1.105,0,2,0.895,2,2l0,0C266.998,406.107,266.103,407.003,264.998,407.003z"
				/>
				<path
					class="fill-app-bulb-base-light"
					d="M264.998,425.003h-21c-1.105,0-2-0.895-2-2l0,0c0-1.105,0.895-2,2-2h21c1.105,0,2,0.895,2,2l0,0C266.998,424.107,266.103,425.003,264.998,425.003z"
				/>
				<path
					class="fill-app-bulb-base-light"
					d="M264.998,443.003h-21c-1.105,0-2-0.895-2-2l0,0c0-1.105,0.895-2,2-2h21c1.105,0,2,0.895,2,2l0,0C266.998,442.107,266.103,443.003,264.998,443.003z"
				/>
				<path
					class="fill-app-bulb-base-light"
					d="M264.998,461.003h-21c-1.105,0-2-0.895-2-2l0,0c0-1.105,0.895-2,2-2h21c1.105,0,2,0.895,2,2l0,0C266.998,460.107,266.103,461.003,264.998,461.003z"
				/>
				<path
					class="fill-app-bulb-base-light"
					d="M274.998,479.003h-21c-1.105,0-2-0.895-2-2l0,0c0-1.105,0.895-2,2-2h21c1.105,0,2,0.895,2,2l0,0C276.998,478.107,276.103,479.003,274.998,479.003z"
				/>
				<path
					class="fill-app-bulb-base-light"
					d="M294.998,496.003h-21c-1.105,0-2-0.895-2-2l0,0c0-1.105,0.895-2,2-2h21c1.105,0,2,0.895,2,2l0,0C296.998,495.107,296.103,496.003,294.998,496.003z"
				/>

				<!-- Right Side Rim Reflection -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-opacity {isOn
						? 'fill-app-bulb-rim opacity-100'
						: 'fill-white/12 opacity-80'}"
					d="M435.998,186.013c0.22,30.93-9.89,59.48-27.07,82.43c-21.04,28.11-32.93,62-32.93,97.11v4.45c0,16.57-13.43,30-30,30h-92c-16.57,0-30-13.43-30-30v-2.96c0-1.53-0.02-3.06-0.08-4.59c5.79,7.36,15.31,15.55,29.08,15.55c24.87,0,48.17-8.37,66.86-22.74c18.68-14.38,32.76-34.76,39.16-58.8c4.46-16.76,11.737-33.687,18.98-50.46c38-88,8-206-135.73-182.17c17.8-8.37,37.72-12.98,58.72-12.83C375.018,51.533,435.468,111.983,435.998,186.013z"
				/>

				<!-- Top-Left Glossy Glass Glare / Specular Highlight Paths -->
				<path
					style="transition-duration: var(--duration);"
					class="transition-opacity {isOn
						? 'fill-app-bulb-glare opacity-100'
						: 'fill-white/30 opacity-90'}"
					d="M242,285c0,0-100.101-55.825-42.511-159.922c1.906-3.446,6.164-4.826,9.724-3.144l22.259,10.511
					c3.875,1.83,5.408,6.53,3.371,10.3C225.474,160.092,206.19,210.197,242,285z"
				/>
				<path
					style="transition-duration: var(--duration);"
					class="transition-opacity {isOn
						? 'fill-app-bulb-glare opacity-100'
						: 'fill-white/45 opacity-95'}"
					d="M215.412,108.706l22.549,11.274c2.292,1.146,5.062,0.543,6.676-1.447
					c3.118-3.847,8.937-10.334,17.103-16.539c3.179-2.416,2.728-7.333-0.82-9.164l-20.567-10.615
					c-1.549-0.799-3.379-0.809-4.939-0.032c-4.599,2.291-14.425,8.047-21.968,18.422
					C211.461,103.334,212.394,107.197,215.412,108.706z"
				/>
			</g>
		</g>
	</svg>
</button>
