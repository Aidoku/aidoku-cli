<script lang="ts">
	export let text: string;
	export let badgeColor: string;
	export let extraClasses = "";
</script>

<span
	class="rounded px-1 py-0.5 text-[0.625rem] leading-3 {extraClasses}"
	style="background-color: {badgeColor};"
>
	{text}
</span>
