<script setup lang="ts">
import { ref } from 'vue';
import router from '../router';

const userId = ref("admin");
const password = ref("password");

if(localStorage.getItem("userId")) {
	router.replace("/dashboard");
}

function login() {
	localStorage.setItem("userId", userId.value);
	router.push({ path: "/dashboard" })
}

</script>

<template>
	<div id="login">
		<div>
			<h1>Login</h1>
			<div>
				<input type="text" v-model="userId" placeholder="user id">
				<input type="password" v-model="password" placeholder="password">
			</div>
			<div>
				<button class="btn btn-nobg clr-green" @click="login">login</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#login{
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	color: white;
	background-color: rgb(40, 40, 40);
}

#login > div{
	padding: 1rem;
	border-radius: 6px;
	background-color: rgb(60, 60, 60);
}

h1{
	font-size: 1.75rem;
	text-align: center;
}

#login > div > *{
	margin-bottom: 1rem;
}

input{
	display: block;
	margin-bottom: 0.5rem;
}

</style>
