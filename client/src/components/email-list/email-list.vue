<script setup lang="ts">
import { ref } from 'vue';
import type { EmailDetails } from '../../types/dbTypes';
import { FilterCycler } from "../../utils/email-filter-state";

const props = defineProps<{
	emailList: EmailDetails[],
	filterers: { [_:string]: FilterCycler }
}>();

const toShowErrIdxs = ref<Set<number>>(new Set);
const emit = defineEmits(["filterChange"])

function toggleErrMsg(idx: number) {
	if(props.emailList[idx].errorMsg.String === "") return;
	if(toShowErrIdxs.value.has(idx)) toShowErrIdxs.value.delete(idx);
	else toShowErrIdxs.value.add(idx);
}

function emitter() {
	emit("filterChange");
}

</script>

<template>
	<div class="email-list scroll-bar">
		<div class="details-title-container">
			<div class="details-title">
				<div @click="filterers.isValidSyntax.cycle(emitter)" :class="filterers.isValidSyntax.className">syntax</div>
				<div @click="filterers.reachable.cycle(emitter)" :class="filterers.reachable.className">reachable</div>
				<div @click="filterers.isDeliverable.cycle(emitter)" :class="filterers.isDeliverable.className">deliverable</div>
				<div @click="filterers.isHostExists.cycle(emitter)" :class="filterers.isHostExists.className">host-exists</div>
				<div @click="filterers.hasMxRecords.cycle(emitter)" :class="filterers.hasMxRecords.className">mx-records</div>
				<div @click="filterers.isDisposable.cycle(emitter)" :class="filterers.isDisposable.className">disposable</div>
				<div @click="filterers.isCatchAll.cycle(emitter)" :class="filterers.isCatchAll.className">catch-all</div>
				<div @click="filterers.isInboxFull.cycle(emitter)" :class="filterers.isInboxFull.className">inbox-full</div>
			</div>
		</div>
		<div 
			v-for="email, idx in props.emailList"
			:key="idx"
			:class="['email', email.errorMsg.String !== '' ? 'error-mail' : '']"
			@click="toggleErrMsg(idx)"
		>
			<div>
				<div>{{ email.emailId }}</div>
				<div v-if="toShowErrIdxs.has(idx)" class="error-msg">{{ email.errorMsg.String }}</div>
			</div>
			<div class="state" v-if="email.errorMsg.Valid">
				<div :class="[email.isValidSyntax ? 'green' : 'red']"></div>
				<div :class="[email.reachable === 'yes' ? 'green' : (email.reachable === 'unknown' ? 'yellow' : 'red')]"></div>
				<div :class="[email.isDeliverable ? 'green' : 'red']"></div>
				<div :class="[email.isHostExists ? 'green' : 'red']"></div>
				<div :class="[email.hasMxRecords ? 'green' : 'red']"></div>
				<div :class="[email.isDisposable ? 'green' : 'red']"></div>
				<div :class="[email.isCatchAll ? 'green' : 'red']"></div>
				<div :class="[email.isInboxFull ? 'green' : 'red']"></div>
			</div>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

.email-list{
	display: flex;
	flex-direction: column;
	position: relative;
	overflow-y: auto;
	height: 100%;
	margin-bottom: 1rem;
	padding: 0.5rem;
	background-color: rgb(40, 40, 40);
	border-radius: 6px;
	width: 100%;
}

.details-title-container{
	display: flex;
	flex-direction: column;
	align-items: flex-end;
	width: 100%;
	position: sticky;
	top: -1rem;
	z-index: 100;
	margin-bottom: 1rem;
	background-color: rgb(40, 40, 40);
	padding: 0.5rem;
	padding-top: 1rem;
}

.details-title{
	display: flex;
	flex-direction: row;
	justify-content: space-around;
	width: 40%;
	max-width: 20rem;
}

.details-title div {
	writing-mode: vertical-rl;
	transform: rotate(180deg);
	text-orientation: sideways-right;
	user-select: none;
	cursor: pointer;
	padding: 0.5rem;
}

.email{
	display: flex;
	flex-direction: row;
	justify-content: space-between;
	width: 100%;
	position: relative;
	border-radius: 6px;
	margin-bottom: 0.2rem;
	background-color: rgb(60, 60, 60);
	padding: 0.5rem;
	font-family: monospace;
}

.state{
	display: flex;
	flex-direction: row;
	justify-content: space-around;
	width: 40%;
	max-width: 20rem;
}

.state div {
	width: 0.8rem;
	height: 0.8rem;
}

.green {
	background-color: rgb(0, 255, 0);
	color: black;
}

.red {
	background-color: rgb(255, 0, 0);
	color: white;
}

.yellow {
	background-color: rgb(255, 255, 0);
	color: black;
}

.error-mail{
	outline: 2px solid rgb(255, 90, 90);
	outline-offset: -2px;
	cursor: pointer;
	transition: 0.2s;
}

.error-mail:hover{
	background-color: rgb(90, 60, 60);
}

.error-msg{
	margin-top: 0.5rem;
	color: rgb(255, 150, 150);
}

</style>
