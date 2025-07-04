export type EmailFile = {
	id: number,
	fileName: string,
}

export type EmailDetails = {
	fileId: number,
	emailId: string,
	isValidSyntax: boolean,
	reachable: string,
	isDeliverable: boolean,
	isHostExists: boolean,
	hasMxRecords: boolean,
	isDisposable: boolean,
	isCatchAll: boolean,
	isInboxFull: boolean,
	errorMsg: { String: string, Valid: boolean },
}

export class FileStats {
	fileId: number = 0
	fileName: string = ""
	createdDateTime: string = ""
	totalEmails: number = 0
	invalidSyntax: number = 0
	reachable: number = 0
	unknown: number = 0
	deliverable: number = 0
	catchAll: number = 0
	disposable: number = 0
	inboxFull: number = 0
	hostExists: number = 0
	errored: number = 0
}

export class ProxyDetails {
	id: number = -1
	userId: string = ""
	proto: string = "socks5"
	host: string = ""
	port: string = ""
	name: string = ""
	password: string = ""
	isInUse: boolean = false
	isEnabled: boolean = true
}
