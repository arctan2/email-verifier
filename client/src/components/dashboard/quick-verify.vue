<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { fetchPost } from '../../utils/fetch';
import EmailList from "../email-list/email-list.vue";
import type { EmailDetails } from '../../types/dbTypes';
import { filterers, getToFilter, resetFilters } from '../../utils/email-filter-state';

enum Msg {
	Verifying = "Verifying...",
	None = ""
}

const emailInput = ref("");
const textAreaElement = ref<null | HTMLTextAreaElement>(null);
const msg = ref<string>(Msg.None);
const results = ref<EmailDetails[]>([]);
const filteredResults = ref<EmailDetails[]>([]);

function resizeTextArea() {
	const target = textAreaElement.value;
	if(!target) return;
	
	target.style.height = "";
	target.style.height = `${Math.min(target.scrollHeight + 4, 16 * 20)}px`;
}

async function verify() {
	const emails = emailInput.value.split("\n").filter(e => e.trim() !== "");
	if(emails.length === 0) return;

	emailInput.value = "";
	resizeTextArea();

	msg.value = Msg.Verifying;

	const data = await fetchPost<{results: EmailDetails[]}>("/verify-emails", emails);

	if(data.err) {
		msg.value = data.msg;
	} else {
		results.value = data.results;
		resizeTextArea();
		msg.value = Msg.None;
	}
	onFilterChange();
}

function cmp(email: EmailDetails, toFilter: string[]) {
	let ret = true;
	for(const k of toFilter) {
		ret = ret && ((email as any)[k] === filterers[k].cmpVal);
	}
	return ret;
}

onMounted(resetFilters);

function onFilterChange() {
	const toFilter = getToFilter();

	if(toFilter.length === 0) {
		filteredResults.value = results.value;
		return;
	}

	filteredResults.value = results.value.filter(email => {
		return cmp(email, toFilter);
	})
}

</script>

<template>
	<div id="quick-verify" class="scroll-bar">
		<div>
			<textarea 
				name="emails-list"
				id="emails-list"
				@input="resizeTextArea"
				v-model="emailInput"
				placeholder="Enter email(s)"
				ref="textAreaElement"
				:disabled="msg === Msg.Verifying"
			></textarea>
			<button class="btn btn-nobg clr-green" @click="verify" :disabled="msg === Msg.Verifying">Verify</button>

			<div v-if="msg !== Msg.None" class="msg">{{ msg }}</div>

			<div id="results" v-if="results.length">
				<EmailList :email-list="filteredResults" :filterers="filterers" @filter-change="onFilterChange" />
			</div>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#quick-verify{
	width: 100%;
	height: 84dvh;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	color: white;
}

#quick-verify > div{
	display: flex;
	flex-direction: column;
	padding: 1rem;
	border-radius: 6px;
	height: 90%;
	width: 96%;
	max-width: 60rem;
	background-color: rgb(60, 60, 60);
}

#emails-list{
	height: calc(2rem + 8px);
	font-size: 1rem;
	padding: 8px;
	margin-bottom: 1rem;
	border-radius: 6px;
	outline: none;
	max-width: 100%;
	resize: none;
	background-color: rgb(50, 50, 50);
	color: white;
}

#emails-list:focus{
	outline: 2px solid #7700ff;
}

button{
	width: max-content;
	margin-bottom: 1rem;
}

.msg{
	text-align: center;
}

#results{
	max-height: 80%;
}

</style>

