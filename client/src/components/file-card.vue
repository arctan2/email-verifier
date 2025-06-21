<script setup lang="ts">
import { FileStats } from "../types/dbTypes";
import { type Percentage } from "../types/fileTypes";
// @ts-ignore
import CircleProgress from "vue3-circle-progress";

const { file: f } = defineProps<{
	file: FileStats,
	percentages: Percentage
}>();

</script>

<template>
	<div class="file-container">
		<h3>{{ f.fileName }}</h3>
		<div class="created-date-time">{{ f.createdDateTime }}</div>
		<slot name="afterFileName"></slot>
		<div>
			<div class="left-section">
				<div class="progress-container">
					<circle-progress class="progress" :percent="percentages.ok.value"
						:fill-color="percentages.ok.color" :empty-color="'rgba(100, 100, 100, 0.3)'" />
					<div class="progress-text" :style="{ color: percentages.ok.color }">
						{{ percentages.ok.value }}% <div>Ok</div>
					</div>
				</div>
				<slot name="afterProgressContainer"></slot>
			</div>
			<div class="stats-container">
				<div class="stats-labels">
					<div>Total:</div>
					<div>Reachable:</div>
					<div>Deliverable:</div>
					<div>Catch-All:</div>
					<div>Unknown:</div>
					<div>Invalid:</div>
					<div>Disposable:</div>
					<div>Error:</div>
				</div>
				<div class="stats-values">
					<div class="stats-num">
						<div>{{ f.totalEmails }}</div>
						<div>{{ f.reachable }}</div>
						<div>{{ f.deliverable }}</div>
						<div>{{ f.catchAll }}</div>
						<div>{{ f.unknown }}</div>
						<div>{{ f.invalidSyntax }}</div>
						<div>{{ f.disposable }}</div>
						<div>{{ f.errored }}</div>
					</div>
					<div class="stats-per">
						<span></span>
						<span class="green">{{ percentages.reachable }}</span>
						<span class="green">{{ percentages.deliverable }}</span>
						<span class="yellow">{{ percentages.catchAll }}</span>
						<span class="yellow">{{ percentages.unknown }}</span>
						<span class="red">{{ percentages.invalidSyntax }}</span>
						<span class="red">{{ percentages.disposable }}</span>
						<span class="red">{{ percentages.errored }}</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

.file-container {
	position: relative;
	padding: 1rem;
	min-height: 21rem;
	min-width: 30rem;
	width: 32%;
	background-color: rgb(60, 60, 60);
	border-radius: 10px;
	box-shadow: 0 5px 15px rgb(0, 0, 0, 0.2);
	font-family: Verdana, Geneva, Tahoma, sans-serif;
}

.file-container h3 {
	font-size: 1.25rem;
}

.created-date-time{
	color: rgb(150, 150, 150);
	font-size: 0.8rem;
	font-family: monospace;
	margin-top: 0.2rem;
	margin-bottom: 1rem;
}

.file-container>div {
	display: flex;
	flex-direction: row;
}

.left-section {
	width: 50%;
	display: flex;
	flex-direction: column;
	align-items: center;
}

.left-section > div:first-child{
	margin-bottom: 1.5rem;
}

.left-section > div button {
	margin-right: 1rem;
}

.progress-container {
	display: flex;
	justify-content: center;
	align-items: center;
	margin-right: 1rem;
}

.progress {
	max-width: 100%;
	max-height: max-content;
}

.progress-text {
	position: absolute;
	text-align: center;
	font-size: 1.5rem;
	font-weight: bold;
}

.progress-text>div {
	font-size: 1rem;
}

.stats-container {
	display: flex;
	flex-direction: row;
	align-items: center;
}

.stats-container>div {
	height: 100%;
	margin-right: 0.5rem;
}

.stats-labels > div, .stats-values > div > * {
	margin-bottom: 0.75rem;
}

.stats-values{
	display: flex;
	flex-direction: row;
}

.stats-values > div{
	display: flex;
	flex-direction: column;
	margin-left: 1rem;
}

.stats-values span {
	display: inline-block;
	border-radius: 3px;
	min-width: 4ch;
}

.stats-values span::after {
	content: "%";
}

.stats-per :first-child::after {
	content: " ";
	white-space: pre;
}

.stats-values .green {
	color: #34eb7a;
	background-color: #34eb7a33;
	outline: 4px solid #34eb7a33;
}

.stats-values .red {
	color: #ff7777;
	background-color: #ff777733;
	outline: 4px solid #ff777733;
}

.stats-values .yellow {
	color: #ffff77;
	background-color: #ffff7733;
	outline: 4px solid #ffff7733;
}

</style>
