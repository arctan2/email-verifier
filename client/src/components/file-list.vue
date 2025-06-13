<script setup lang="ts">
import { ref, shallowRef } from 'vue';
import type { EmailFile } from '../types/dbTypes';
import { fetchGet, fetchPostWithHeader, fetchDelete } from '../utils/fetch';
import { setPopupError, setPopupInfo } from "../utils/popup-data";
import DeleteIcon from '../assets/icons/delete.svg';
import Modal from '../common-components/modal.vue';

const props = defineProps<{
	setCurSelected: (f: EmailFile) => void,
	curSelected: EmailFile | null
}>();

const filesList = shallowRef<EmailFile[]>([]);
const uploadingList = shallowRef<EmailFile[]>([]);
const msg = ref<string>("Loading...");
const deleteIdx = ref<number>(-1);
const editIdx = ref<number>(-1);
const showLoadingMessage = ref<string>("");

async function fetchAllAccounts() {
	const data = await fetchGet<{ allFiles: EmailFile[] }>("/get-all-files");
	if (data.err) {
		msg.value = data.msg;
	} else if (data.allFiles.length === 0) {
		msg.value = "No files found.";
	} else {
		filesList.value = data.allFiles;
		msg.value = "";
	}
}

async function uploadFile(f: File) {
	const formData = new FormData();
	formData.append("file", f)

	uploadingList.value = [...uploadingList.value, { id: -1, fileName: f.name }];

	const data = await fetchPostWithHeader<{ id: number, fileName: string }>(
		"/upload-file", {}, formData
	);

	uploadingList.value = uploadingList.value.filter(a => a.fileName !== data.fileName);

	if(data.err) {
		setPopupError(data.msg);
	} else {
		const newFiles = [...filesList.value, { id: data.id, fileName: data.fileName }];
		newFiles.sort((a, b) => a.id - b.id);
		filesList.value = newFiles;
	}

	if(filesList.value.length > 0) {
		msg.value = "";
	}
}

async function deleteTransaction() {
	showLoadingMessage.value = "Deleting...";
	const id = filesList.value[deleteIdx.value].id;
	deleteIdx.value = -1;
	const res = await fetchDelete<{keywords: string[]}>(`/delete-file?id=${id}`);
	if(res.err) {
		setPopupError(res.msg);
	} else {
		setPopupInfo("File deleted successfully.");
	}
	showLoadingMessage.value = "";

	filesList.value = [];
	fetchAllAccounts();
}

function handleFileSelect(e: Event) {
	const input = e.target as HTMLInputElement;
	const filesAsArray = Array.from(input?.files || []);

	for(const f of filesAsArray) {
		uploadFile(f);
	}
}

fetchAllAccounts();

</script>

<template>
	<div id="files-list-container">
		<Modal v-if="deleteIdx >= 0" :dismiss="() => deleteIdx = -1">
			<div class="modal delete-modal">
				<h2 style="text-align: center;">Delete</h2>

				<div>Are you sure you want to delete the file "{{ filesList[deleteIdx].fileName }}"?</div>

				<div class="btn-container" style="align-self: center;">
					<button class="btn btn-nobg clr-red" @click="deleteTransaction">delete</button>
					<button class="btn btn-nobg" @click="deleteIdx = -1">cancel</button>
				</div>
			</div>
		</Modal>

		<h1>Files</h1>
		<div id="files-list" class="scroll-bar">
			<div v-if="msg !== '' && uploadingList.length === 0" style="text-align: center;">{{ msg }}</div>
			<template v-else>
				<div v-for="f, idx in filesList" :key="idx">
					<div :class="['file', props.curSelected?.id === f.id ? 'selected-file' : '']" @click.self="props.setCurSelected(f)">
						<div class="file-name">{{ f.fileName }}</div>
						<div class="icon-btns" v-if="editIdx < 0">
							<div @click="deleteIdx = idx">
								<img :src="DeleteIcon" alt="">
							</div>
						</div>
					</div>
				</div>

				<div v-for="f, idx in uploadingList" :key="idx">
					<div class="file uploading">
						<div class="file-name">{{ f.fileName }}</div>
					</div>
				</div>
			</template>

			<div>
				<label for="upload-excel-files" class="upload-input">Click to upload files</label>
				<input type="file" accept=".txt" id="upload-excel-files" @change="handleFileSelect" multiple>
			</div>
			<div class="refreshing" v-if="showLoadingMessage !== ''">{{ showLoadingMessage }}</div>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#files-list-container {
	padding: 0.5rem;
	background-color: rgb(255, 255, 255);
	box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3);
	border-radius: 6px;
	width: 100%;
	max-width: 20rem;
	height: 100%;
}

#files-list {
	display: flex;
	flex-direction: column;
	height: 89dvh;
	overflow-y: auto;
	overflow-x: hidden;
	padding: 0.5rem;
	border: 2px solid grey;
	border-radius: 6px;
}

h1 {
	margin-bottom: 1rem;
}

.file {
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: space-between;
	border-radius: 6px;
	border: 2px solid rgb(72, 0, 255);
	margin-bottom: 0.8rem;
	padding: 0.5rem;
	cursor: pointer;
	transition: all 0.2s;
	font-family: monospace;
	position: relative;
}

.selected-file {
	background-color: rgba(255, 255, 0, 0.5);
}

.file.selected-file:hover{
	background-color: rgba(255, 255, 0, 0.5);
}

.file>div {
	margin-bottom: 0.3rem;
}

.file:hover {
	background-color: rgba(72, 0, 255, 0.2);
}

.file-id {
	font-size: 1.1rem;
	font-weight: bold;
	color: #5632a8;
}

.file-name {
	color: #333333;
	pointer-events: none;
	user-select: none;
}

input[type=file] {
	display: none;
}

.upload-input {
	display: flex;
	justify-content: center;
	align-items: center;
	width: 100%;
	height: 4rem;
	border: 3px dashed #1111ff;
	border-radius: 6px;
	font-family: monospace;
	cursor: pointer;
	font-size: 1rem;
	opacity: 0.5;
	transition: 0.2s;
	user-select: none;
}

.upload-input:hover {
	opacity: 1;
}

.uploading::after{
	content: "";
	width: 1.5rem;
	height: 1.5rem;
	position: absolute;
	right: 1rem;
	border-left: 3px solid black;
	border-radius: 10rem;
	animation: spin 1s infinite linear;
}

@keyframes spin {
	0%{
		transform: rotate(0deg);
	}
	100%{
		transform: rotate(360deg);
	}
}

.icon-btns{
	display: flex;
	flex-direction: row;
	align-items: center;
}

.icon-btns > div {
	--size: 2rem;
	width: var(--size);
	height: var(--size);
	margin: 0 0.3rem;
	padding: 0.2rem;
	border-radius: 6px;
	cursor: pointer;
	user-select: none;
	transition: 0.2s;
}

.icon-btns > div:hover{
	box-shadow: 0 0 10px rgba(255, 0, 0, 0.5);
	background-color: rgba(255, 0, 0, 0.2);
}

.icon-btns > div img {
	width: 100%;
	height: 100%;
	object-fit: contain;
	pointer-events: none;
}

.refreshing{
	position: absolute;
	display: flex;
	align-items: center;
	justify-content: center;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.7);
	border-radius: 6px;
	color: white;
}

.delete-modal{
	justify-content: space-between;
	background-color: rgb(255, 255, 255);
	border: 2px solid #ff6161;
}

</style>
