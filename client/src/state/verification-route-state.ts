import { ref, readonly } from "vue";

interface VerificationProps {
	toVerifyCount: number,
	fileId: number,
}

const data = ref<VerificationProps>({
	toVerifyCount: -1,
	fileId: -1
});

export function setVerificationProps(d: VerificationProps) {
	data.value = d;
}

export const verificationProps = readonly(data);
