/// <reference path="alpine.d.ts" />

(function (root) {
  const namesInEnglish = new Intl.DisplayNames(["en"], { type: "language" });

  /**
   * @param {string} code
   * @returns {string}
   */
  function simpleLanguageName(code) {
    if (code === "multi") {
      return "Multi-Language";
    } else {
      return namesInEnglish.of(code);
    }
  }

  /**
   * @param {string} code
   * @returns {string}
   */
  function languageName(code) {
    if (code === "multi") {
      return "Multi-Language";
    }

    if (code === "en") {
      return "English";
    }

    const namesInNative = new Intl.DisplayNames([code], { type: "language" });
    return `${namesInEnglish.of(code)} - ${namesInNative.of(code)}`;
  }

  const LoadingStatus = {
    Loading: "loading",
    Loaded: "loaded",
    Error: "error",
  };

  const NsfwOption = {
    All: "all",
    Safe: "safe",
    Nsfw: "nsfw",
  };

  document.addEventListener("alpine:init", () => {
    Alpine.store(
      "sourceUrl",
      root.location.href.replace(root.location.hash, "")
    );
    Alpine.store(
      "addUrl",
      `aidoku://addSourceList?url=${root.location.href.replace(
        root.location.hash,
        ""
      )}`
    );

    Alpine.bind("TriStateCheckbox", () => ({
      indeterminate: true,

      '@click'() {
        const event = this.$event;
        event.preventDefault();

        if (event.target.indeterminate) {
          event.target.checked = true;
          event.target.indeterminate = false;
        } else if (event.target.checked) {
          event.target.checked = false;
        } else if (!event.target.checked) {
          event.target.indeterminate = true;
        }

        this.$dispatch("sl-change");
      }
    }))

    Alpine.data("sourceList", () => ({
      LoadingStatus,
      NsfwOption,

      simpleLanguageName,
      languageName,
      /**
       * @type {Source[]}
       */
      sources: [],

      /**
       * @type {string[]}
       */
      languages: [],

      /**
       * @type {"loading" | "loaded" | "error"}
       */
      loading: LoadingStatus.Loading,

      // options
      /**
       * @type {Source[]}
       */
      filtered: [],
      
      query: "",

      /**
       * @type {string[]}
       */
      selectedLanguages: [],

      nsfw: NsfwOption.All,

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

            const langLhs = simpleLanguageName(lhs.lang);
            const langRhs = simpleLanguageName(rhs.lang);
            return (
              langLhs.localeCompare(langRhs) || lhs.name.localeCompare(rhs.name)
            );
          });
          this.languages = [
            ...new Set(this.sources.map((source) => source.lang)),
          ];
          this.loading = LoadingStatus.Loaded;
        } catch {
          this.loading = LoadingStatus.Error;
        }

        if (this.filtered.length === 0) {
          this.updateFilteredList();
        }
        this.$nextTick(() => {
          root.location.hash && root.location.replace(root.location.hash);
        });
      },

      updateFilteredList() {
        this.filtered = this.sources
          .filter((item) =>
            this.query
              ? item.name.toLowerCase().includes(this.query.toLowerCase()) ||
                item.id.toLowerCase().includes(this.query.toLowerCase())
              : true
          )
          .filter((item) => 
            this.nsfw === NsfwOption.All 
              ? true 
              : this.nsfw === NsfwOption.Nsfw
              ? item.nsfw > 1
              : item.nsfw < 2)
          .filter((item) =>
            this.selectedLanguages.length
              ? this.selectedLanguages.includes(item.lang)
              : true
          );
      },
    }));
  });
})(typeof self !== "undefined" ? self : this);
