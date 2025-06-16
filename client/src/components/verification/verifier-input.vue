<script setup lang="ts">
import { ref } from 'vue';
import { setPopupError } from '../utils/popup-data';

const props = defineProps<{
	toVerifyCount: number,
	createVerifier: (batchSize: number, retryCount: number, delayMs: number, proxies: string[]) => void
}>();

const batchSizeInput = ref(1000);
const delayMsInput = ref(1000);
const retryCountInput = ref(3);
const proxiesInput = ref("");

function createVerifier() {
	let batchSize = Number(batchSizeInput.value);
	if(Number.isNaN(batchSize) || batchSize <= 0) {
		setPopupError("Invalid batch size");
		return;
	}
	batchSize = Math.floor(batchSize);

	let retryCount = Number(retryCountInput.value);
	if(Number.isNaN(retryCount) || retryCount <= 0) {
		setPopupError("Invalid retry count");
		return;
	}
	retryCount = Math.floor(retryCount);

	let delayMs = Number(delayMsInput.value);
	if(Number.isNaN(delayMs) || delayMs <= 0) {
		setPopupError("Invalid delay(ms)");
		return;
	}
	delayMs = Math.floor(delayMs);

	props.createVerifier(
		batchSize,
		retryCount,
		delayMs,
		proxiesInput.value.split("\n").filter(e => e !== '')
	)
}

</script>

<template>
	<div class="verify-inputs">
		<h1>Create verifier</h1>
		<div>To verify emails: {{ props.toVerifyCount }}</div>
		<div>Batch size: <input type="number" v-model="batchSizeInput" /></div>
		<div>Retry count: <input type="number" v-model="retryCountInput" /></div>
		<div>Delay(ms): <input type="number" v-model="delayMsInput" /></div>
		<div>
			<div>Proxies:</div>
			<textarea name="proxies" id="proxies" v-model="proxiesInput"></textarea>
		</div>
		<div>
			<button class="btn btn-nobg clr-lblue" @click="createVerifier">create</button>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

.verify{
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

.verify-inputs{
	display: flex;
	flex-direction: column;
}

.verify-inputs > div{
	margin-bottom: 0.5rem;
}

h1{
	font-size: 1.5rem;
	margin-bottom: 1rem;
}

#proxies{
	width: 25rem;
	padding: 0.5rem;
	resize: none;
}
</style>
