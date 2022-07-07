<script lang="ts">
	import SourceItem from "./SourceItem.svelte";
	import type { Source } from "../types/source";
	export let sources: Source[];
</script>

<div class="source-list">
	{#each sources as source}
		<SourceItem {...source} />
		<hr />
	{/each}
</div>
