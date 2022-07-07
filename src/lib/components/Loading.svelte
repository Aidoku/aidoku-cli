<div aria-label="Loading" class="spinner relative inline-block h-14 w-14">
	<div class="bar0" />
	<div class="bar1" />
	<div class="bar2" />
	<div class="bar3" />
	<div class="bar4" />
	<div class="bar5" />
	<div class="bar6" />
	<div class="bar7" />
</div>

<style lang="scss">
	.spinner > div {
		@apply absolute left-[49%] top-[43%] h-[16%] w-[6%] animate-fade rounded-[50px] bg-[#78787a] opacity-0 shadow-[0_0_3px_rgba(0,0,0,0.2)];

		@for $i from 0 through 7 {
			&.bar#{$i} {
				transform: rotate(45deg * $i) translate(0, -130%);
				animation-delay: 0s + 0.125s * $i;
			}
		}
	}
</style>
