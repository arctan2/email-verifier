<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { Socket } from "../../socket/socket"
import { API_URL } from '../../utils/fetch';
import VerifierInput from "./verifier-input.vue";
import VerifierDetailsComp from "./verifier-details.vue";
import VerifierProgress from "./verifier-progress.vue";
import VerifierDone from "./verifier-done.vue";
import { setPopupError } from '../../utils/popup-data';
import { Status, type VerifierDetails } from "../../types/verifierTypes";
import { verificationProps } from '../../state/verification-route-state';
import router from '../../router';

const props = verificationProps.value;
const ws = new Socket();
const emit = defineEmits(['verification-complete'])

const msg = ref("Loading...");
const curStatus = ref<Status>(Status.NotCreated);
const verifierDetails = ref<null | VerifierDetails>(null);
const errMsg = ref<string>("");

function createVerifier(batchSize: number, retryCount: number, delayMs: number, proxies: string[]) {
	const details = {
		emailCount: props.toVerifyCount,
		batchSize,
		retryCount,
		delayMs,
		proxies
	}

	ws.emit("create-verifier", details);
}

function collectVerifierInfo() {
	ws.emit("get-verifier-details");
	msg.value = "Collecting verifier info...";
}

function listenWs() {
	ws.on("status", (status: Status) => {
		curStatus.value = status;
		if(status !== Status.NotCreated) {
			collectVerifierInfo();
		}
		msg.value = "";
	})

	ws.on("create-verifier-res", (res: { err: boolean, msg: string }) => {
		if(res.err) {
			setPopupError(res.msg);
		} else {
			curStatus.value = Status.Created;
			collectVerifierInfo();
		}
	})

	ws.on("get-verifier-details-res", (details: VerifierDetails) => {
		if((details as any).err) {
			msg.value = (details as any).msg;
			return;
		}
		if(curStatus.value === Status.Running && details.state === Status.Done) {
			emit("verification-complete");
		}
		verifierDetails.value = details;
		curStatus.value = details.state;
		msg.value = "";
	})

	ws.on("run-verifier-err", (err: string) => {
		curStatus.value = Status.Err;
		errMsg.value = err;
	})
}

function runVerifier() {
	ws.emit("run-verifier");
}

function removeVerifier() {
	ws.emit("remove-verifier");
	curStatus.value = Status.NotCreated;
	verifierDetails.value = null;
}

function onDoneClick() {
	removeVerifier();
	router.replace({ name: "Emails", query: { fileId: props.fileId } });
}

onMounted(async () => {
	listenWs();
	await ws.connect(API_URL(`/api/web/${props.fileId}/verification-ws`), _ => {
		msg.value = "Error in making websocket connection";
	});

	if(ws.isConnected) {
		msg.value = "Connected. Waiting for status...";
	}
})

onUnmounted(() => {
	ws.close();
})

</script>

<template>
	<div class="main-section">
		<div v-if="msg !== ''">{{ msg }}</div>
		<div v-else class="verify-container">
			<VerifierInput
				v-if="curStatus === Status.NotCreated"
				:to-verify-count="props.toVerifyCount"
				:create-verifier="createVerifier"
			/>

			<VerifierDetailsComp 
				v-if="curStatus === Status.Created && verifierDetails !== null"
				:verifier-details="verifierDetails"
				:run-verifier="runVerifier"
				:cancel-verifier="removeVerifier"
			/>

			<VerifierProgress
				v-if="curStatus === Status.Running && verifierDetails !== null"
				:verifier-details="verifierDetails"
				:ws="ws"
			/>

			<VerifierDone
				v-if="curStatus === Status.Done && verifierDetails !== null"
				:verifier-details="verifierDetails"
				:done="onDoneClick"
			/>

			<div v-if="curStatus === Status.Err" style="color: red">{{ errMsg }}</div>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

.main-section {
	height: 100%;
	overflow: hidden;
}

.verify-container{
	width: 100%;
	height: 100%;
	background-color: rgb(50, 50, 50);
	border-radius: 6px;
	position: relative;
	padding: 0.5rem;
	color: white;
}

</style>
