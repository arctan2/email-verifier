<script setup lang="ts">
import { computed, ref } from 'vue';

const props = defineProps<{
	length: number,
	toShow: number,
	onChange: (newPage: number) => void
}>()

const curPage = ref<number>(0);

const toShow = computed(() => {
	let toShow = props.toShow;

	if((props.length < (toShow * 2)) && (props.length > toShow)) {
		toShow = Math.max(Math.floor(toShow / 2), props.length - toShow);
	}

	return toShow;
})

const state = computed(() => {
	const showLast = curPage.value < (props.length - toShow.value);
	const showFirst = curPage.value >= toShow.value;
	let start = 0;
	let select = curPage.value + 1;

	if(showFirst && showLast) {
		select = Math.ceil(toShow.value / 2);
		start = curPage.value - select + 1;
	}

	if(!showLast && showFirst) {
		start = props.length - toShow.value;
		select = toShow.value + 1 - (props.length - curPage.value)
	}

	return {
		start,
		select,
		showLast,
		showFirst 
	}
})

function setCurPage(pageNum: number) {
	const prev = curPage.value;
	if(pageNum < 0) pageNum = 0;
	if(pageNum >= props.length) pageNum = props.length - 1;
	curPage.value = pageNum;
	if(prev !== pageNum) {
		props.onChange(pageNum);
	}
}

</script>

<template>
<div class="page-nav">
	<button class="prev" @click="setCurPage(curPage - 1)"><</button>

	<button v-if="state.showFirst" @click="setCurPage(0)">1</button>

	<span v-if="state.showFirst" style="display: flex; align-items: center;">...</span>

	<button
		v-for="i in Math.min(toShow, props.length)"
		@click="setCurPage(state.start + i - 1)"
		:class="i === state.select ? 'cur-page' : ''"
	>{{state.start + i}}</button>

	<span v-if="state.showLast" style="display: flex; align-items: center;">...</span>

	<button v-if="state.showLast" @click="setCurPage(props.length - 1)">{{ props.length }}</button>

	<button class="next" @click="setCurPage(curPage + 1)">></button>
</div>
</template>

<style scoped>
@import "@css/common.css";

.page-nav{
	width: 100%;
	height: 2rem;
	display: flex;
	flex-direction: row;
	justify-content: center;
	user-select: none;
}

.page-nav > button{
	font-family: monospace;
	border: none;
	outline: none;
	background-color: transparent;
	color: white;
	font-size: 1rem;
	min-width: 2rem;
	min-height: 2rem;
	border-radius: 6px;
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	transition: all 0.1s;
}

.page-nav > button:hover{
	background-color: rgb(255, 255, 255, 0.2);
}

.page-nav > button:active{
	background-color: rgb(255, 255, 255, 0.3);
}

.page-nav > .cur-page:hover{
	background-color: rgb(0, 170, 255);
}

.page-nav > .cur-page{
	background-color: rgb(0, 170, 255);
	color: black;
}

.page-nav > .next, .page-nav > .prev{
	border: 1px solid grey;
	margin: 0 0.5rem;
	border-radius: 50rem;
}

</style>
