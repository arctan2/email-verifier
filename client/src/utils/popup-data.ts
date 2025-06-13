import { reactive } from "vue";

type PopupData = { 
	type: "info" | "error" | "",
	text: string 
}

export const popupData = reactive<PopupData>({
	type: "",
	text: ""
})

const popupStack: PopupData[] = [];

export function setPopupError(msg: string, stack: boolean = true) {
	if(stack && popupData.type !== "") {
		popupStack.push({ type: popupData.type, text: popupData.text })
	}
	popupData.type = "error";
	popupData.text = msg;
}

export function setPopupInfo(msg: string, stack: boolean = true) {
	if(stack && popupData.type !== "") {
		popupStack.push({ type: popupData.type, text: popupData.text })
	}
	popupData.type = "info";
	popupData.text = msg;
}

export function hidePopup() {
	const data = popupStack.pop();
	if(data) {
		if(data.type === "error") setPopupError(data.text, false);
		else setPopupInfo(data.text, false);
		return;
	}
	popupData.type = "";
}

