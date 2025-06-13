export type EmailFile = {
	id: number,
	fileName: string,
}

export type EmailDetails = {
	fileId: number,
	emailId: string,
	isValidSyntax: boolean,
	isReachable: boolean,
	isDeliverable: boolean,
	isHostExists: boolean,
	hasMxRecords: boolean,
	isDisposable: boolean,
	isCatchAll: boolean,
	isInboxFull: boolean,
	errorMsg: string,
}

