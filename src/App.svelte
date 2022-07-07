<script lang="ts">
	import SourcesWrapper from "./lib/components/SourcesWrapper.svelte";

	const baseUrl = window.location.href.replace(window.location.hash, "");
</script>

<header
	class="mx-auto my-4 flex max-w-2xl flex-col items-center justify-center text-center text-xl"
>
	<p class="my-2 pb-2 font-semibold text-aidoku-color">An Aidoku source list</p>
	<a
		class="align-center inline-block justify-center rounded-md bg-light-tint px-6 py-2 text-base font-medium text-aidoku-color hover:bg-button-gradient"
		href="aidoku://addSourceList?url={baseUrl}"
	>
		Add to Aidoku
	</a>
</header>
<blockquote class="mx-auto my-3 w-11/12 max-w-2xl">
	On a device with Aidoku installed, tap <a
		class="font-normal text-aidoku-color"
		href="aidoku://addSourceList?url={baseUrl}">Add to Aidoku</a
	> or use the base URL to add this source list.
</blockquote>
<div
	class="mx-auto my-3 w-11/12 max-w-2xl rounded-sm border border-blue-400 bg-blue-50 py-3 px-4"
>
	<p class="mb-2 font-semibold">Base URL:</p>
	<p>{baseUrl}</p>
</div>
<div class="mx-auto my-3 w-11/12 max-w-2xl">
	<h1 class="pb-4 text-3xl font-semibold">Sources</h1>
	<SourcesWrapper />
</div>
