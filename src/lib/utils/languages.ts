const namesInEnglish = new Intl.DisplayNames(["en"], { type: "language" });

export function languageName(code: string): string[] {
	if (code == "multi") {
		return ["Multi-Language", ""];
	} else {
		const namesInNative = new Intl.DisplayNames([code], { type: "language" });

		// All the language codes are probably valid if they got merged.
		// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
		return [namesInEnglish.of(code)!, namesInNative.of(code)!];
	}
}

export function fullLanguageName(code: string): string {
	const names = languageName(code);
	return names[1]
		? code !== "en"
			? `${names[0]} - ${names[1]}`
			: names[0]
		: names[0];
}
