<script setup lang="ts">
import type { VerifierDetails } from '../types/verifierTypes';
import { computeBatchesList } from './verifier';

const props = defineProps<{
	verifierDetails: VerifierDetails,
	done: () => void
}>();

const verifierDetails = props.verifierDetails;

const allBatchList = computeBatchesList(
	verifierDetails.batchSize,
	verifierDetails.emailCount,
	-1,
	verifierDetails.completedBatches
);

</script>

<template>
	<div class="verifier-done">
		<h1>Verification Complete</h1>
		<div class="verifier-progress scroll-bar">
			<div class="batches">
				<div
					v-for="batch in allBatchList"
					:key="batch.num"
					:class="['batch', ...batch.classNames]"
				>
					<div class="batch-title">Batch {{ batch.from === batch.to ? batch.from : `${batch.from}-${batch.to}` }}</div>
					<div class="batch-body scroll-bar">
						<div v-for="r, i in batch.retries" class="retry-item" :key="i">
							<span style="color: #93b7f5">Try-{{i}}:</span>
							<span class="success-result">{{r.success}}/{{r.total}}</span>
							<div class="fail-result">{{r.failed}} failed</div>
							<div class="fail-result">{{r.retry}} retry</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<button class="btn btn-nobg clr-orange" @click="props.done">Done</button>
	</div>
</template>

<style scoped>
@import "@css/common.css";

.verifier-done{
	display: flex;
	flex-direction: column;
	position: relative;
	overflow-y: auto;
	height: 100%;
	margin-bottom: 1rem;
	padding: 0.5rem;
	border-radius: 6px;
	width: 100%;
}

button{
	max-width: max-content;
	align-self: center;
}

.verifier-progress{
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

.batches{
	display: flex;
	flex-direction: row;
	flex-wrap: wrap;
	gap: 1rem;
}

.batch{
	min-width: 10rem;
	min-height: 4rem;
	background-color: rgb(64, 64, 64);
	border-radius: 6px;
	padding: 0.5rem;
	border: 2px solid #808080;
}

.current-batch{
	border: 2px solid #00ff1a;
}

.not-started-batch{
	opacity: 0.2;
}

.batch-title{
	color: yellow;
	font-weight: bold;
	margin-bottom: 0.5rem;
	position: relative;
}

.current-batch .batch-title::after{
	content: "";
	min-width: 1rem;
	min-height: 1rem;
	position: absolute;
	right: 0;
	margin-left: 0.5rem;
	border-left: 3px solid rgb(0, 255, 0);
	border-radius: 10rem;
	animation: spin 1s infinite linear;
}

.batch-body{
	--h: 10rem;
	min-height: var(--h);
	max-height: var(--h);
	background-color: rgb(30, 30, 30);
	border-radius: 6px;
	padding: 0.5rem;
	overflow-y: auto;
	font-family: monospace;
}

.retry-item, .progress{
	margin-bottom: 0.5rem;
}

.retry-item{
	padding: 0.5rem;
	border-radius: 6px;
	background-color: rgb(60, 60, 60);
}

.success-result{
	color: rgb(0, 255, 0);
}

.fail-result{
	color: rgb(250, 110, 110);
}

</style>
