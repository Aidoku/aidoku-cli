<script lang="ts">
	import { onMount, tick } from "svelte";
	import { languageName } from "../utils/languages";
	import SourceList from "./SourceList.svelte";
	import Loading from "./Loading.svelte";
	import SourceFilter from "./SourceFilter.svelte";
	import type { Source } from "../types/source";

	enum LoadingMode {
		Loading,
		Loaded,
		Error,
	}

	let sources: Source[] = [];
	let filteredSources: Source[] = [];
	let loading = LoadingMode.Loading;

	onMount(async () => {
		try {
			const res = await fetch(`./index.min.json`);
			sources = (await res.json()).sort((lhs: Source, rhs: Source) => {
				if (lhs.lang === "multi" && rhs.lang !== "multi") {
					return -1;
				}
				if (lhs.lang !== "multi" && rhs.lang === "multi") {
					return 1;
				}
				if (lhs.lang === "en" && rhs.lang !== "en") {
					return -1;
				}
				if (rhs.lang === "en" && lhs.lang !== "en") {
					return 1;
				}

				const [langLhs] = languageName(lhs.lang);
				const [langRhs] = languageName(rhs.lang);
				if (langLhs < langRhs) {
					return -1;
				}
				if (langRhs > langLhs) {
					return 1;
				}
				return lhs.name < rhs.name ? -1 : lhs.name === rhs.name ? 0 : 1;
			});
			loading = LoadingMode.Loaded;
		} catch {
			loading = LoadingMode.Error;
		}
	});

	$: tick().then(() => {
		if (window.location.hash && loading === LoadingMode.Loaded) {
			const elem = document.getElementById(
				window.location.hash.replace("#", "")
			) as HTMLElement;
			elem.scrollIntoView();
		}
	});
</script>

{#if loading === LoadingMode.Loading}
	<div class="flex items-center justify-center">
		<Loading />
	</div>
{:else if loading === LoadingMode.Error}
	<div class="flex items-center justify-center">
		<div class="text-red-500">
			<p>Could not load sources.</p>
			<p>This source list might not work in the app.</p>
		</div>
	</div>
{:else}
	<SourceFilter {sources} bind:filtered={filteredSources} />
	<SourceList sources={filteredSources} />
{/if}
