<script setup lang="ts">
import type { VerifierDetails } from '../types/verifierTypes';

const props = defineProps<{
	verifierDetails: VerifierDetails,
	runVerifier: () => void
}>();

console.log(props.verifierDetails)

</script>

<template>
	<div class="verifier-details">
		<h1>Verifier</h1>
		<div>
			<button class="btn btn-nobg clr-green" @click="props.runVerifier">Run</button>
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
