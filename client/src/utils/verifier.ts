import { FileStats } from "../types/dbTypes";
import { type Percentage } from "../types/fileTypes";

export function calcPercentages(f: FileStats) {
	const total = f.totalEmails;
	const reachable = Math.floor(f.reachable / total * 100);
	const deliverable = Math.floor(f.deliverable / total * 100);
	const ok = Math.max(reachable, deliverable);
	let okColor = "#34eb7a";

	if (ok <= 35) {
		okColor = "#ff7777";
	} else if (ok <= 65) {
		okColor = "#ffff77";
	}

	const p: Percentage = {
		ok: { color: okColor, value: ok },
		reachable,
		deliverable,
		catchAll: Math.floor(f.catchAll / total * 100),
		invalidSyntax: Math.floor(f.invalidSyntax / total * 100),
		unknown: Math.floor(f.unknown / total * 100),
		disposable: Math.floor(f.disposable / total * 100),
		errored: Math.floor(f.errored / total * 100),
	};

	return p;
}
