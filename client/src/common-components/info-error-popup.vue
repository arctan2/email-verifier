<script setup lang="ts">
import { hidePopup, popupData } from '../utils/popup-data';

</script>

<template>
	<div v-if="popupData.type === 'info'" id="popup-container" class="info" @click.self="hidePopup">
		<div>
			<div>
				{{popupData.text}}
			</div>
			<button class="btn btn-nobg" @click="hidePopup">ok</button>
		</div>
	</div>
	<div v-if="popupData.type === 'error'" id="popup-container" class="error" @click.self="hidePopup">
		<div>
			<div>
				{{popupData.text}}
			</div>
			<button class="btn btn-nobg" @click="hidePopup">ok</button>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#popup-container{
	display: flex;
	justify-content: center;
	align-items: center;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.6);
	position: absolute;
	top: 0;
	left: 0;
	z-index: 100;
}

#popup-container > div{
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	min-width: 15rem;
	min-height: 8rem;
	background-color: rgb(255, 255, 255);
	color: black;
	border-radius: 4px;
	font-size: 1.2rem;
	box-shadow: 0 8px 30px rgb(0, 0, 0);
	padding: 1rem;
}

#popup-container.error > div {
	color: rgb(255, 57, 43);
	border: 2px solid rgb(255, 57, 43);
}

#popup-container button{
	margin-top: 1.5rem;
}

#info-popup{
	max-width: 40rem;
	min-width: 15rem;
	min-height: 3rem;
	background-color: rgb(40, 40, 40);
	border: 2px solid rgb(255, 200, 100);
	border-radius: 4px;
	box-shadow: 0 8px 30px rgba(0, 0, 0, 0.5);
	position: absolute;
	bottom: 15%;
	left: 50%;
	transform: translateX(-50%);
	padding: 0.5rem 1rem;
	display: flex;
	flex-direction: row;
	align-items: center;
	z-index: 100;
}

#info-popup .icon{
	display: inline-flex;
	justify-content: center;
	align-items: center;
	min-width: 1.2rem;
	min-height: 1.2rem;
	width: 1.2rem;
	height: 1.2rem;
	background-color: rgb(255, 255, 0);
	color: black;
	font-family: monospace;
	border-radius: 20px;
	margin-right: 0.5rem;
}

</style>../utils/popup-data
