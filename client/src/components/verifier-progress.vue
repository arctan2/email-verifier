<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { Socket } from '../socket/socket';
import type { VerifierDetails } from '../types/verifierTypes';
import type { ProgressData } from '../types/verifierTypes';
import { type Batch, computeBatchesList } from "./verifier";

const { verifierDetails, ws } = defineProps<{
	verifierDetails: VerifierDetails,
	ws: Socket
}>();

const allBatchList = ref<Batch[]>([]);
const curBatch = ref<number>(-1);
const curBatchProgress = ref<ProgressData[]>([])
const curProgress = ref<ProgressData>({
	total: 0,
	progress: 0,
	failed: 0,
	retry: 0,
	success: 0
})
const curBatchBody = ref<HTMLDivElement[]>([]);
const isDelay = ref(false);

function listenWs() {
	ws.on("batch-start", (data: string) => {
		isDelay.value = false;
		const batchNum = Number(data);
		const b = allBatchList.value[batchNum];

		if(curBatch.value !== -1) {
			const prevBatch = allBatchList.value[curBatch.value];
			prevBatch.retries = [...curBatchProgress.value];
			prevBatch.classNames = [];
		}

		curBatch.value = batchNum;
		curProgress.value = {
			success: 0,
			failed: 0,
			retry: 0,
			progress: 0,
			total: b.to - b.from + 1,
		};
		curBatchProgress.value = [];
		b.classNames = ["current-batch"];
	})

	ws.on("progress", (data: ProgressData) => {
		curProgress.value = data;
	})

	ws.on("retry-delay", (data: ProgressData) => {
		curBatchProgress.value.push(data);

		setTimeout(() => {
			const e = curBatchBody.value.pop();

			if(e) {
				e.scrollTop = e.scrollHeight;
			}
		}, 10)

		isDelay.value = true;
	})

	ws.on("after-all-retries", (data: ProgressData) => {
		curBatchProgress.value.push(data);
		setTimeout(() => {
			const e = curBatchBody.value.pop();

			if(e) {
				e.scrollTop = e.scrollHeight;
			}
		}, 10)
		isDelay.value = false;
	})

	ws.on("batch-delay", () => {
		isDelay.value = true;
	})

	ws.on("retry-begin", (data: ProgressData) => {
		isDelay.value = false;
		curProgress.value = data;
	})
}

onMounted(() => {
	curBatch.value = verifierDetails.currentBatchNumber;
	curBatchProgress.value = verifierDetails.currentProgressList;
	curBatchProgress.value.pop();

	allBatchList.value = computeBatchesList(
		verifierDetails.batchSize,
		verifierDetails.emailCount,
		curBatch.value,
		verifierDetails.completedBatches
	);
	listenWs();
})

</script>

<template>
	<div class="verifier-progress scroll-bar">
		<div class="batches">
			<div
				v-for="batch in allBatchList"
				:key="batch.num"
				:class="['batch', ...batch.classNames]"
			>
				<div class="batch-title">Batch {{ batch.from === batch.to ? batch.from : `${batch.from}-${batch.to}` }}</div>
				<div class="batch-body scroll-bar" :ref="curBatch === batch.num ? 'curBatchBody' : undefined">
					<template v-if="batch.num !== curBatch">
						<div v-for="r, i in batch.retries" class="retry-item" :key="i">
							<span style="color: #93b7f5">Try-{{i}}:</span>
							<span class="success-result">{{r.success}}/{{r.total}}</span>
							<div class="fail-result">{{r.failed}} failed</div>
							<div class="fail-result">{{r.retry}} retry</div>
						</div>
					</template>

					<template v-else-if="curBatchProgress.length > 0">
						<div v-for="r, i in curBatchProgress" class="retry-item" :key="i">
							<span style="color: #93b7f5">Try-{{i}}:</span>
							<span class="success-result">{{r.success}}/{{r.total}}</span>
							<div class="fail-result">{{r.failed}} failed</div>
							<div class="fail-result">{{r.retry}} retry</div>
						</div>
					</template>

					<template v-if="batch.num === curBatch">
						<div v-if="isDelay" class="delay">Delay</div>
						<div class="progress" v-else>
							Progress: <span class="retry-result">
								{{curProgress.progress}}/{{curProgress.total}}
							</span>
						</div>
					</template>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

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
	display: inline-block;
	min-width: 1rem;
	min-height: 1rem;
	position: relative;
	top: 0.1rem;
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

.retry-result{
	color: #93b7f5;
}

.delay{
	position: relative;
	margin-bottom: 0.5rem;
	color: yellow;
}

.delay::after{
	content: "";
	display: inline-block;
	min-width: 1rem;
	min-height: 1rem;
	position: relative;
	top: 0.1rem;
	margin-left: 0.5rem;
	border-left: 3px solid yellow;
	border-radius: 10rem;
	animation: spin 1s infinite linear;
}

.progress::after{
	content: "";
	display: inline-block;
	position: relative;
	min-width: 1rem;
	min-height: 1rem;
	margin-left: 0.5rem;
	top: 0.1rem;
	border-left: 3px solid white;
	border-radius: 10rem;
	animation: spin 1s infinite linear;
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

@keyframes spin {
	0%{
		transform: rotate(0deg);
	}
	100%{
		transform: rotate(360deg);
	}
}

</style>
