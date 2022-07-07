<script lang="ts">
	import { languageName } from "../utils/languages";
	import Badge from "./Badge.svelte";

	export let id: string;
	export let name: string;
	export let file: string;
	export let icon: string;
	export let lang: string;
	export let version = 1;
	export let nsfw = 0;

	let shortLangName = languageName(lang)[0];

	// These properties are here so a Source object can just be spread.
	// svelte-ignore unused-export-let
	export let minAppVersion = "0.1.0";
	// svelte-ignore unused-export-let
	export let maxAppVersion = "0.1.0";
</script>

<div class="align-center group flex px-2 py-2 lg:px-6" {id}>
	<a
		href="#{id}"
		class="source-anchor float-left my-auto -ml-6 px-2 text-xl font-normal text-aidoku-color opacity-0 hover:underline group-hover:opacity-100"
		>#</a
	>
	<div class="my-auto">
		<img
			class="mr-2 rounded-lg"
			alt="Icon for {name}"
			src="./icons/{icon}"
			loading="lazy"
			width="42"
			height="42"
		/>
	</div>
	<div class="flex flex-1 flex-col">
		<div class="flex flex-row flex-wrap items-center gap-1 font-semibold">
			{name}
			<span class="text-sm text-gray-400">v{version}</span>
			{#if nsfw == 2}
				<Badge text="18+" badgeColor="#ff3b304d" extraClasses="" />
			{/if}
		</div>
		<div class="text-sm text-gray-400">
			{shortLangName}
		</div>
	</div>
	<div class="mx-auto my-auto">
		<button
			class="rounded-full bg-gray-100 py-1 px-4 text-xs font-bold text-aidoku-color no-underline hover:bg-button-gradient"
		>
			<a href="./sources/{file}" download>DOWNLOAD</a>
		</button>
	</div>
</div>
