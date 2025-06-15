<script setup lang="ts">
import FileList from "../components/file-list.vue";
import FileDetails from "../components/file-details.vue";
import { ref } from "vue";
import type { EmailFile } from "../types/dbTypes";

const curSelectedFile = ref<null | EmailFile>(null);

</script>

<template>
	<div id="files">
		<FileList
			:curSelected="curSelectedFile"
			:setCurSelected="(f: null | EmailFile) => curSelectedFile = f?.id === curSelectedFile?.id ? null : f"
		/>
		<FileDetails 
			:curSelected="curSelectedFile"
		/>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#files {
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	position: relative;
	background-color: rgb(40, 40, 40);
	width: 100%;
	height: 100%;
	padding: 0.2rem;
}

@media (max-width: 500px) {
	#files {
		padding: 1%;
	}
}

</style>
