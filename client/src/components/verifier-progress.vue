<script setup lang="ts">
import { onMounted } from 'vue';
import { Socket } from '../socket/socket';
import type { VerifierDetails } from '../types/verifierTypes';

const props = defineProps<{
	verifierDetails: VerifierDetails,
	ws: Socket
}>();

console.log(props.verifierDetails)

function listenWs() {
	const ws = props.ws;

	ws.on("progress", (data: any) => {
		console.log("progress", data);
	})

	ws.on("batch-start", (data: any) => {
		console.log("batch-start", data);
	})

	ws.on("delay", (data: any) => {
		console.log("delay", data);
	})

	ws.on("done", (data: any) => {
		console.log("done", data);
	})

	ws.on("after-all-retries", (data: any) => {
		console.log("after-all-retries", data);
	})

	ws.on("retry-begin", (data: any) => {
		console.log("retry-begin", data);
	})
}

onMounted(() => {
	listenWs();
})

</script>

<template>
	<div class="verifier-details">
		<h1>Verifier</h1>
		<div>
			Progress
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

.verifier-details{
	display: flex;
	flex-direction: column;
	position: relative;
	overflow-y: auto;
	height: 92%;
	margin-bottom: 1rem;
	padding: 0.5rem;
	background-color: rgb(40, 40, 40);
	border-radius: 6px;
	width: 100%;
}

h1{
	font-size: 1.5rem;
	margin-bottom: 1rem;
}

</style>
