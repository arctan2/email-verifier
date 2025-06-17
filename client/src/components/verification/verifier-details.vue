<script setup lang="ts">
import type { VerifierDetails } from '../../types/verifierTypes';

const props = defineProps<{
	verifierDetails: VerifierDetails,
	runVerifier: () => void,
	cancelVerifier: () => void
}>();

</script>

<template>
	<div class="verifier-details-container">
		<h1>Verifier</h1>
		<div class="verifier-details">
			<div>To verify: {{ verifierDetails.emailCount }}</div>
			<div>Batch size: {{ verifierDetails.batchSize }}</div>
			<div>Delay(ms): {{ verifierDetails.delayMs }}</div>
			<div>Retry count: {{ verifierDetails.retryCount }}</div>
			<div>Proxies: <div v-for="p in verifierDetails.proxies">{{ p }}</div></div>

			<div>
				<button class="btn btn-nobg clr-green" @click="props.runVerifier">run</button>
				<button class="btn btn-nobg clr-white" @click="cancelVerifier">cancel</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

.verifier-details-container{
	display: flex;
	flex-direction: column;
	position: relative;
	overflow-y: auto;
	height: 92%;
	margin-bottom: 1rem;
	padding: 0.5rem;
	border-radius: 6px;
	width: 100%;
}

h1{
	font-size: 1.5rem;
	margin-bottom: 1rem;
}

.verifier-details > div{
	margin-bottom: 0.5rem;
}

button {
	margin-top: 1rem;
	margin-right: 1rem;
}

</style>
