<script setup lang="ts">
const props = defineProps<{
	min: number,
	max: number,
	step: number,
	dir: "rtl" | "ltr",
	value: number
}>();
</script>

<template>
<div class="range">
	<input
		type="range"
		:value="props.value"
		:min="props.min"
		:max="props.max" 
		:step="props.step"
		:class="props.dir"
		@input="$emit('input', $event)" 
	/>
</div>
</template>

<style scoped>
*{
	--control-bg: rgb(26, 0, 61);
	--purple: rgb(112, 0, 255);
	--aqua: rgb(3, 252, 161);
	--pink: rgb(255, 0, 64);
	--element-height: 1.5rem;
}

.range{
	display: flex;
	justify-content: center;
	align-items: center;
}

.rtl{
	direction: rtl;	
}

.ltr{
	direction: ltr;	
}

.range input[type="range"]{
	--track-height: 0.2rem;
	--thumb-width: 1rem;
	--thumb-height: 0.5rem;
	-webkit-appearance: none;
	appearance: none;
	cursor: pointer;
	background: transparent;
}

.range input[type="range"]:hover{
	--thumb-width: 1rem;
	--thumb-height: 1rem;
}

.range input[type="range"]::-webkit-slider-runnable-track{
	--bg: var(--aqua);
	background: var(--bg);
	height: var(--track-height);
	border-radius: 1rem;
	transition: all 0.3s;
}

.range input[type="range"]:hover::-webkit-slider-runnable-track{
	--track-height: 0.4rem;
	box-shadow: 0 0 5px var(--bg);
}

.range input[type="range"]::-webkit-slider-thumb{
	--bg: white;
	-webkit-appearance: none;
	appearance: none;
	width: var(--thumb-width);
	height: var(--thumb-height);
	background: white;
	box-shadow: 0 0 1px var(--bg);
	border-radius: 1rem;
	margin-top: calc(var(--thumb-height) / 2 * -1 + var(--track-height) / 2);
	transition: all 0.3s;
}

.range input[type="range"]:hover::-webkit-slider-thumb{
	box-shadow: 0 0 8px var(--bg);
}
</style>
		
