<script setup lang="ts">
import { watch, ref } from 'vue';
import { fetchGet } from '../utils/fetch';
import type { EmailFile } from '../types/dbTypes';
import EmailsDetailsList from "./email-details-list.vue";
import Verification from "./verification/verification.vue";

const props = defineProps<{
	curSelected: null | EmailFile
}>();

enum Tab {
	List = "List",
	Verification = "Verification"
}

const msg = ref<string>("Loading...");
const emailsCount = ref(0);
const toVerifyCount = ref(0);
const curTab = ref<Tab>(Tab.List);

async function fetchFileDetails(fileId: number) {
	const data = await fetchGet<{ emailsCount: number, toVerifyCount: number }>(`/${fileId}/get-file-details`);
	if (data.err) {
		msg.value = data.msg;
	} else {
		emailsCount.value = data.emailsCount;
		toVerifyCount.value = data.toVerifyCount;
		msg.value = "";
	}
}

function goToList() {
	curTab.value = Tab.List;
	fetchFileDetails(props.curSelected?.id as number);
}

watch(props, (cur) => {
	if(cur.curSelected !== null) {
		msg.value = "Loading...";
		curTab.value = Tab.List;
		fetchFileDetails(cur.curSelected.id);
	}
})

</script>

<template>
	<div id="emails-container">
		<div v-if="props.curSelected === null"></div>
		<div v-else-if="msg !== ''" style="text-align: center;">{{ msg }}</div>
		<template v-else>
			<div class="file-details" v-if="props.curSelected">
				<div>file name: {{ props.curSelected.fileName }}</div>
				<div>number of emails: {{ emailsCount }}</div>
				<div class="button-container">
					<button 
						v-for="t in Tab"
						:class="['btn', 'btn-nobg', curTab === t ? 'selected-button' : '']" @click="curTab = t"
					>{{t}}</button>
				</div>
			</div>
			<EmailsDetailsList
				v-if="curTab === Tab.List"
				:totalEmailCount="emailsCount"
				:fileId="props.curSelected.id"
			/>
			<Verification
				v-if="curTab === Tab.Verification"
				:to-verify-count="toVerifyCount"
				:curFile="props.curSelected"
				:go-to-list="goToList"
			/>
		</template>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#emails-container {
	display: flex;
	flex-direction: column;
	position: relative;
	width: 100%;
	height: 100%;
	padding: 0.5rem;
	border-radius: 6px;
	margin-left: 0.2rem;
	background-color: rgb(60, 60, 60);
	color: white;
}

.file-details{
	display: flex;
	flex-direction: column;
	padding: 1rem;
	border-radius: 6px;
	box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
	margin-bottom: 0.5rem;
	border: 2px solid #666666;
}

.button-container {
	margin-top: 0.5rem;
	display: flex;
	flex-direction: row;
}

.button-container button {
	margin-right: 0.5rem;
}

.selected-button{
	background-color: white;
	color: black;
}

</style>
