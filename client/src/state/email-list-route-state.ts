import { ref, readonly } from "vue";

interface EmailListProps {
	totalEmailCount: number,
	fileId: number
}

const data = ref<EmailListProps>({
	totalEmailCount: -1,
	fileId: -1
});

export function setEmailListProps(d: EmailListProps) {
	data.value = d;
}

export const emailListProps = readonly(data);
