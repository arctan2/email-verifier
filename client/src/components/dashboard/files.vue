<script setup lang="ts">
import { shallowRef, ref, computed } from 'vue';
import { FileStats, type EmailFile } from '../../types/dbTypes';
import { type Percentage } from '../../types/fileTypes';
import { fetchDelete, fetchGet, fetchPostWithHeader } from '../../utils/fetch';
import { calcPercentages } from '../../utils/verifier';
import { setPopupError, setPopupInfo } from '../../utils/popup-data';
import DeleteIcon from "../../assets/icons/delete.svg";
import Modal from '../../common-components/modal.vue';
import router from '../../router';
import FileCard from "../file-card.vue";

const filesList = shallowRef<FileStats[]>([]);
const msg = ref<string>("Loading...");
const showLoadingMessage = ref<string>("");
const uploadingList = shallowRef<EmailFile[]>([]);
const deleteIdx = ref<number>(-1);

async function fetchFilesStatsList() {
	const data = await fetchGet<{ statsList: FileStats[] }>(`/get-file-list-stats?userId=${localStorage.getItem("userId")}`);
	if (data.err) {
		msg.value = data.msg;
	} else {
		filesList.value = data.statsList;
		msg.value = "";
	}
}

fetchFilesStatsList();

const percentages = computed(() => {
	const percentageList: Percentage[] = []

	for(const f of filesList.value) {
		percentageList.push(calcPercentages(f));
	}

	return percentageList;
});

async function uploadFile(f: File) {
	const formData = new FormData();
	formData.append("file", f)

	uploadingList.value = [...uploadingList.value, { id: -1, fileName: f.name }];

	const data = await fetchPostWithHeader<{ id: number, fileName: string, emailCount: number }>(
		"/upload-file", {}, formData
	);

	uploadingList.value = uploadingList.value.filter(a => a.fileName !== f.name);

	if(data.err) {
		setPopupError(data.msg);
	} else {
		const newFile = new FileStats();
		newFile.fileId = data.id;
		newFile.fileName = data.fileName;
		newFile.totalEmails = data.emailCount;

		const newFiles: FileStats[] = [...filesList.value, newFile];
		newFiles.sort((a, b) => b.fileId - a.fileId);
		filesList.value = newFiles;
	}

	if(filesList.value.length > 0) {
		msg.value = "";
	}
}

async function deleteFile() {
	showLoadingMessage.value = "Deleting...";
	const id = filesList.value[deleteIdx.value].fileId;
	deleteIdx.value = -1;
	const res = await fetchDelete<{keywords: string[]}>(`/delete-file?id=${id}`);
	if(res.err) {
		setPopupError(res.msg);
	} else {
		setPopupInfo("File deleted successfully.");
	}
	showLoadingMessage.value = "";

	filesList.value = [];
	fetchFilesStatsList();
}

function handleFileSelect(e: Event) {
	const input = e.target as HTMLInputElement;
	const filesAsArray = Array.from(input?.files || []);

	for(const f of filesAsArray) {
		uploadFile(f);
	}
	input.value = "";
}

function navigateVerify(file: FileStats) {
	router.push({ path: "/file/verification", query: { fileId: file.fileId } });
}

function navigateDetails(file: FileStats) {
	router.push({ path: "/file/emails", query: { fileId: file.fileId } });
}

</script>

<template>
	<div id="files">
		<Modal v-if="deleteIdx >= 0" :dismiss="() => deleteIdx = -1">
			<div class="modal delete-modal">
				<h2 style="text-align: center;">Delete</h2>

				<div>Are you sure you want to delete the file "{{ filesList[deleteIdx].fileName }}"?</div>

				<div class="btn-container" style="align-self: center;">
					<button class="btn btn-nobg clr-red" @click="deleteFile">delete</button>
					<button class="btn btn-nobg" @click="deleteIdx = -1">cancel</button>
				</div>
			</div>
		</Modal>

		<div v-if="msg !== '' && uploadingList.length === 0" style="text-align: center; margin-bottom: 1rem;">{{ msg }}</div>
		<template v-else>
			<div id="files-list" class="scroll-bar">
				<div class="upload-file">
					<label for="upload-excel-files" class="upload-input">Click to upload files</label>
					<input type="file" accept=".txt" id="upload-excel-files" @change="handleFileSelect" multiple>
				</div>

				<div v-for="f, idx in uploadingList" :key="idx" class="file-container">
					<h3>{{ f.fileName }}</h3>
					<div class="uploading">
						uploading 
					</div>
				</div>

				<FileCard v-for="f, i in filesList" :key="i" :percentages="percentages[i]" :file="f">
					<template v-slot:afterFileName>
						<div class="delete-icon" @click="deleteIdx = i">
							<img :src="DeleteIcon" alt="">
						</div>
					</template>

					<template v-slot:afterProgressContainer>
						<div>
							<button class="btn btn-nobg" @click="navigateDetails(f)">Details</button>
							<button class="btn btn-nobg clr-green" @click="navigateVerify(f)">Verify</button>
						</div>
					</template>
				</FileCard>
			</div>
		</template>
		<div class="refreshing" v-if="showLoadingMessage !== ''">{{ showLoadingMessage }}</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#files {
	width: 100%;
	height: 100%;
	color: white;
	position: relative;
}

#files-list {
	display: flex;
	flex-direction: row;
	flex-wrap: wrap;
	gap: 1rem;
	overflow-y: auto;
	overflow-x: hidden;
	padding: 0.5rem;
	max-height: 84dvh;
}

.refreshing {
	position: absolute;
	display: flex;
	align-items: center;
	justify-content: center;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.7);
	z-index: 100;
	border-radius: 6px;
	color: white;
}

.upload-input {
	display: flex;
	justify-content: center;
	align-items: center;
	border: 3px dashed rgb(3, 148, 252);
	padding: 1rem;
	font-family: monospace;
	cursor: pointer;
	font-size: 1rem;
	opacity: 0.5;
	transition: 0.2s;
	user-select: none;
	position: relative;
	padding: 1rem;
	height: 20rem;
	min-width: 30rem;
	width: 32%;
	background-color: rgb(60, 60, 60);
	border-radius: 10px;
	box-shadow: 0 5px 15px rgb(0, 0, 0, 0.2);
}

.file-container {
	position: relative;
	padding: 1rem;
	height: 20rem;
	min-width: 30rem;
	width: 32%;
	background-color: rgb(60, 60, 60);
	border-radius: 10px;
	box-shadow: 0 5px 15px rgb(0, 0, 0, 0.2);
	font-family: Verdana, Geneva, Tahoma, sans-serif;
}

.file-container h3 {
	font-size: 1.25rem;
	margin-bottom: 1rem;
}

.file-container>div {
	display: flex;
	flex-direction: row;
}

.upload-input:hover {
	opacity: 1;
}

.uploading{
	width: max-content;
	position: relative;
}

.uploading::after{
	content: "";
	width: 1.5rem;
	height: 1.5rem;
	position: absolute;
	right: -2rem;
	top: -0.2rem;
	border-left: 3px solid white;
	border-radius: 10rem;
	animation: spin 1s infinite linear;
}

.delete-icon {
	--size: 2rem;
	width: var(--size);
	height: var(--size);
	margin: 0 0.3rem;
	padding: 0.2rem;
	border-radius: 6px;
	cursor: pointer;
	user-select: none;
	transition: 0.2s;
	position: absolute;
	right: 1rem;
	top: 0.5rem;
}

.delete-icon:hover{
	box-shadow: 0 0 10px rgba(255, 0, 0, 0.5);
	background-color: rgba(255, 0, 0, 0.2);
}

.delete-icon img {
	width: 100%;
	height: 100%;
	object-fit: contain;
	pointer-events: none;
}

input[type=file] {
	display: none;
}

.delete-modal{
	justify-content: space-between;
	background-color: rgb(30, 30, 30);
	color: white;
	border: 2px solid #ff6161;
}

.left-section > div button {
	margin-right: 1rem;
}

</style>
