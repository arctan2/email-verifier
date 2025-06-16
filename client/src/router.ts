import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";

import RootRoute from "./root-route.vue";
import Dashboard from "./views/Dashboard.vue";
import DashboardQuickVerify from "./components/dashboard/quick-verify.vue";
import DashboardFiles from "./components/dashboard/files.vue";

const routes: Array<RouteRecordRaw> = [
	{ path: "/", component: RootRoute },
	{ 
		path: "/dashboard", component: Dashboard, 
		children: [
			{
				path: "quick-verify",
				component: DashboardQuickVerify
			},
			{
				path: "files",
				component: DashboardFiles
			},
		]
	},
];

const router = createRouter({
	history: createWebHistory("/"),
	routes,
});

router.beforeEach(to => {
	if(!localStorage.getItem("userId") && to.path !== "/") {
		return { path: "/" }
	}
});


export default router;
