import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";

import RootRoute from "./root-route.vue";
import Dashboard from "./views/Dashboard.vue";
import DashboardQuickVerify from "./components/dashboard/quick-verify.vue";
import DashboardFiles from "./components/dashboard/files.vue";
import FileOfId from "./views/FileOfId.vue";
import EmailDetailsList from "./components/email-list/email-details-list.vue";
import Verification from "./components/verification/verification.vue";

const routes: Array<RouteRecordRaw> = [
	{ path: "/", component: RootRoute },
	{ 
		path: "/dashboard", component: Dashboard, name: 'Dashboard',
		children: [
			{
				name: "QuickVerify",
				path: "quick-verify",
				component: DashboardQuickVerify
			},
			{
				name: "Files",
				path: "files",
				component: DashboardFiles
			},
		]
	},
	{
		path: "/file", component: FileOfId,
		children: [
			{
				name: 'Emails',
				path: "emails",
				component: EmailDetailsList
			},
			{
				name: 'Verification',
				path: "verification",
				component: Verification
			},
		]
	}
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
