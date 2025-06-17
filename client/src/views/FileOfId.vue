<script setup lang="ts">
import { useRoute, RouterView, RouterLink } from 'vue-router';
import router from '../router';
import { fetchGet } from '../utils/fetch';
import { ref } from 'vue';
import { setVerificationProps } from '../state/verification-route-state';
import { setEmailListProps } from '../state/email-list-route-state';
import FileCard from "../components/file-card.vue";
import { FileStats } from '../types/dbTypes';
import { calcPercentages } from '../utils/verifier';
import { setPopupError } from '../utils/popup-data';

const route = useRoute();
const fileId = Number(route.query.fileId);
const msg = ref("Loading...");
const fileStats = ref<FileStats | null>(null);

if(Number.isNaN(fileId)) router.replace("/dashboard");

async function fetchFileDetails() {
	const data = await fetchGet<{ emailsCount: number, toVerifyCount: number }>(`/${fileId}/get-file-details`);
	if (data.err) {
		return data.msg;
	} else {
		setEmailListProps({ fileId: fileId, totalEmailCount: data.emailsCount });
		setVerificationProps({ toVerifyCount: data.toVerifyCount, fileId });
	}
	return "";
}

async function fetchFileStats() {
	const data = await fetchGet<{ fileStats: FileStats }>(
		`/get-file-stats?userId=${localStorage.getItem("userId")}&fileId=${fileId}`
	);
	if (data.err) {
		return data.msg;
	} else {
		fileStats.value = data.fileStats;
	}
	return "";
}

async function fetchDetails() {
	let err = await fetchFileStats();
	if(err !== "") {
		msg.value = err;
		return;
	}
	msg.value = await fetchFileDetails();
}

fetchDetails();

async function verificationComplete() {
	let err = await fetchFileStats();
	if(err !== "") {
		setPopupError(err);
	}

	err = await fetchFileDetails();
	if(err !== "") {
		setPopupError(err);
	}
}

</script>

<template>
	<div v-if="msg !== ''" class="msg">{{ msg }}</div>
	<div id="file-details">
		<div class="left-panel">
			<div v-if="fileStats !== null" class="stats">
				<FileCard :file="fileStats" :percentages="calcPercentages(fileStats)" />
			</div>
			<nav>
				<RouterLink
					:class="route.name === 'Emails' ? 'current' : ''"
					:to="{ path: `/file/emails`, query: { fileId } }"
				>Emails</RouterLink>

				<RouterLink
					:class="route.name === 'Verification' ? 'current' : ''"
					:to="{ path: `/file/verification`, query: { fileId } }"
				>Verification</RouterLink>

				<RouterLink :to="{ path: `/dashboard`, query: { fileId } }">Go Back</RouterLink>
			</nav>
		</div>
		<div class="content" v-if="msg === ''">
			<RouterView @verification-complete="verificationComplete"/>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#file-details{
	width: 100%;
	height: 100%;
	display: flex;
	flex-direction: row;
	color: white;
	background-color: rgb(40, 40, 40);
}

#file-details > div{
	padding: 0.5rem;
	border-radius: 6px;
	background-color: rgb(30, 30, 30);
}

.msg{
	color: white;
	background-color: rgb(30, 30, 30);
	width: 100%;
	height: 100%;
	display: flex;
	justify-content: center;
	align-items: center;
	font-size: 2rem;
}

.left-panel{
	height: 100%;
	min-width: max-content;
	width: 32rem;
	margin-right: 0.3rem;
}

.content{
	width: 100%;
	height: 100%;
}

nav{
	display: flex;
	flex-direction: column;
	margin-top: 1rem;
}

nav > *{
	text-decoration: none;
	color: white;
	border-radius: 6px;
	font-size: 1.2rem;
	margin-bottom: 0.5rem;
	padding: 1rem;
}

nav > .current{
	background-color: white;
	color: black;
	transition: 0.2s;
}

</style>
