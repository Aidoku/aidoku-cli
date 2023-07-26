(function (scope) {
    const namesInEnglish = new Intl.DisplayNames(["en"], { type: "language" });

    /**
     * 
     * @param {string} code 
     * @returns {[string, string]}
     */
    function languageName(code) {
        if (code == "multi") {
            return ["Multi-Language", ""];
        } else {
            const namesInNative = new Intl.DisplayNames([code], { type: "language" });

            // All the language codes are probably valid if they got merged.
            // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
            return [namesInEnglish.of(code), namesInNative.of(code)];
        }
    }

    /**
     * 
     * @param {string} code 
     * @returns {string}
     */
    function fullLanguageName(code) {
        const names = languageName(code);
        return names[1]
            ? code !== "en"
                ? `${names[0]} - ${names[1]}`
                : names[0]
            : names[0];
    }

    const LoadingStatus = {
        Loading: "loading",
        Loaded: "loaded",
        Error: "error",
    };

    document.addEventListener("alpine:init", () => {
        Alpine.store("sourceUrl", window.location.href.replace(window.location.hash, ""))
        Alpine.store("addUrl", `aidoku://addSourceList?url=${window.location.href.replace(window.location.hash, "")}`)

        Alpine.data("sourceList", () => ({
            sources: [],
            languages: [],
            loading: LoadingStatus.Loading,

            LoadingStatus,
            languageName,
            fullLanguageName,

            // options
            filtered: [],
            query: "",
            selectedLanguages: [],
            nsfw: true,

            async init() {
                try {
                    const res = await fetch(`./index.min.json`);
                    this.sources = (await res.json()).sort((lhs, rhs) => {
                        if (lhs.lang === "multi" && rhs.lang !== "multi") {
                            return -1;
                        }
                        if (lhs.lang !== "multi" && rhs.lang === "multi") {
                            return 1;
                        }
                        if (lhs.lang === "en" && rhs.lang !== "en") {
                            return -1;
                        }
                        if (rhs.lang === "en" && lhs.lang !== "en") {
                            return 1;
                        }
        
                        const [langLhs] = languageName(lhs.lang);
                        const [langRhs] = languageName(rhs.lang);
                        return langLhs.localeCompare(langRhs) || lhs.name.localeCompare(rhs.name);
                    });
                    this.languages = [...new Set(this.sources.map((source) => source.lang))];
                    this.loading = LoadingStatus.Loaded;
                } catch {
                    this.loading = LoadingStatus.Error;
                }

                if (scope.location.hash) {
                    this.$nextTick(() => { 
                        scope.location.replace(scope.location.hash);
                    });
                }
            },

            updateFilteredList() {
                this.filtered = this.sources
                    .filter((item) =>
                        this.query
                            ? item.name.toLowerCase().includes(this.query.toLowerCase()) ||
                            item.id.toLowerCase().includes(this.query.toLowerCase())
                            : true
                    )
                    .filter((item) => (this.nsfw ? true : (item.nsfw ?? 0) <= 1))
                    .filter((item) =>
                        this.selectedLanguages.length ? this.selectedLanguages.includes(item.lang) : true
                    );
            }
        }))
    })
})(window);
