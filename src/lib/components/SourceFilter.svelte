<script lang="ts">
	import MultiSelect from "svelte-multiselect";
	import { fullLanguageName } from "../utils/languages";
	import type { Source } from "../types/source";

	export let sources: Source[];
	export let filtered: Source[];

	const languages: string[] = [...new Set(sources.map((item) => item.lang))];
	let query = "";
	let languageQuery: string[] = [];
	let nsfw = true;

	$: filtered = sources
		.filter((item) =>
			query
				? item.name.toLowerCase().includes(query.toLowerCase()) ||
				  item.id.toLowerCase().includes(query.toLowerCase())
				: true
		)
		.filter((item) => (nsfw ? true : (item.nsfw ?? 0) <= 1))
		.filter((item) =>
			languageQuery.length ? languageQuery.includes(item.lang) : true
		);
</script>

<div class="pb-2">
	<div id="search">
		<input
			type="search"
			class="box-border flex w-full cursor-text rounded-sm border border-gray-300 px-2 py-1 leading-6 outline-none placeholder:text-sm placeholder:text-gray-300 focus-within:border-blue-400"
			bind:value={query}
			placeholder="Search by name or ID..."
		/>
	</div>
	<div id="languages">
		<MultiSelect
			bind:selectedValues={languageQuery}
			options={languages.map((language) => ({
				label: fullLanguageName(language),
				value: language,
			}))}
			placeholder="Show specific languages..."
			inputClass="placeholder:text-gray-300 placeholder:text-sm border-gray-300 focus-within:border-blue-400 outline-none"
			--sms-border-radius="0.125rem"
			--sms-padding="0.305rem 0.5rem"
			--sms-placeholder-color="lightgray"
		/>
	</div>
	<div id="nsfw">
		<label>
			Show NSFW sources?
			<input type="checkbox" bind:checked={nsfw} />
		</label>
	</div>
</div>
