<script setup lang="ts">
const props = defineProps<{
	onChange: (s:string) => void,
	value: string,
	options: {[_:string]:string} | Array<string>,
	styleClass?: Array<string>
}>();

</script>

<template>
	<select class="custom-select" @input="(e) => props.onChange((e.target as any).value)">
		<option :selected="alg === props.value" v-for="alg in props.options">{{alg}}</option>
	</select>
</template>

<style scoped>
.custom-select{
	appearance: none;
	-webkit-appearance: none;
	outline: 0;
	width: 7.5rem;
	height: 1.7rem;
	border: 2px solid rgb(255, 255, 255);
	background-color: transparent;
	border-radius: 4px;
	color: white;
	padding-left: 0.5rem;
	font-weight: 600;
	display: block;
	cursor: pointer;
	font-family: Arial, Helvetica, sans-serif;
}

.custom-select > option {
	font-weight: 600;
	background-color: black;
}

</style>
		
